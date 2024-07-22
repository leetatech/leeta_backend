package application

import (
	"context"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/mailer"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/user/domain"
	"go.uber.org/zap"
	"time"
)

type userAppHandler struct {
	tokenHandler  pkg.TokenHandler
	encryptor     pkg.EncryptorManager
	idGenerator   pkg.IDGenerator
	otpGenerator  pkg.OtpGenerator
	logger        *zap.Logger
	EmailClient   mailer.MailerClient
	allRepository pkg.Repositories
}

type UserApplication interface {
	VendorVerification(ctx context.Context, request domain.VendorVerificationRequest) (*pkg.DefaultResponse, error)
	AddVendorByAdmin(ctx context.Context, request domain.VendorVerificationRequest) (*pkg.DefaultResponse, error)
	UpdateUserRecord(ctx context.Context, request models.User) (*pkg.DefaultResponse, error)
	GetAuthenticatedUser(ctx context.Context) (*models.Customer, error)
}

func NewUserApplication(request pkg.DefaultApplicationRequest) UserApplication {
	return &userAppHandler{
		tokenHandler:  request.TokenHandler,
		encryptor:     pkg.NewEncryptor(),
		idGenerator:   pkg.NewIDGenerator(),
		otpGenerator:  pkg.NewOTPGenerator(),
		logger:        request.Logger,
		EmailClient:   request.EmailClient,
		allRepository: request.AllRepository,
	}
}

func (u userAppHandler) VendorVerification(ctx context.Context, request domain.VendorVerificationRequest) (*pkg.DefaultResponse, error) {
	claims, err := u.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	if claims.Role != models.VendorCategory {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	vendorUpdateRequest := domain.VendorDetailsUpdateRequest{
		ID:        claims.UserID,
		Identity:  request.Identity,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Status:    models.Registered,
	}
	err = u.allRepository.UserRepository.VendorDetailsUpdate(vendorUpdateRequest)
	if err != nil {
		return nil, err
	}

	category, err := models.SetBusinessCategory(request.Category)
	if err != nil {
		return nil, err
	}

	business := models.Business{
		ID:          u.idGenerator.Generate(),
		VendorID:    claims.UserID,
		Name:        request.Name,
		CAC:         request.CAC,
		Category:    category,
		Description: request.Description,
		Phone:       request.Phone,
		Address:     request.Address,
		Status:      models.Registered,
		Timestamp:   time.Now().Unix(),
	}
	err = u.allRepository.UserRepository.RegisterVendorBusiness(business)
	if err != nil {
		return nil, err
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Business successfully registered"}, nil
}

func (u userAppHandler) AddVendorByAdmin(ctx context.Context, request domain.VendorVerificationRequest) (*pkg.DefaultResponse, error) {
	claims, err := u.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}
	if claims.Role != models.AdminCategory {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}
	_, err = u.allRepository.AuthRepository.GetAdminByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}

	vendor := models.Vendor{
		User: models.User{
			ID:        u.idGenerator.Generate(),
			FirstName: request.FirstName,
			LastName:  request.LastName,
			Status:    models.Registered,
		},
		AdminID: claims.UserID,
		TimeStamps: models.TimeStamps{
			Time: time.Now().Unix(),
		},
	}
	err = u.allRepository.AuthRepository.CreateUser(ctx, vendor)
	if err != nil {
		return nil, err
	}

	category, err := models.SetBusinessCategory(request.Category)
	if err != nil {
		return nil, err
	}

	business := models.Business{
		ID:          u.idGenerator.Generate(),
		VendorID:    vendor.ID,
		Name:        request.Name,
		CAC:         request.CAC,
		Category:    category,
		Description: request.Description,
		Phone:       request.Phone,
		Address:     request.Address,
		Status:      models.Registered,
		Timestamp:   time.Now().Unix(),
	}
	err = u.allRepository.UserRepository.RegisterVendorBusiness(business)
	if err != nil {
		return nil, err
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Business successfully registered"}, nil
}

func (u userAppHandler) UpdateUserRecord(ctx context.Context, request models.User) (*pkg.DefaultResponse, error) {
	claims, err := u.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	customer, err := u.allRepository.UserRepository.GetCustomerByID(claims.UserID)
	if err != nil {
		return nil, err
	}
	if len(request.Address) > 0 {
		customer.Address = append(customer.Address, request.Address...)
	}

	if request.FirstName != "" {
		customer.FirstName = request.FirstName
	}

	if request.LastName != "" {
		customer.LastName = request.LastName
	}

	err = u.allRepository.UserRepository.UpdateUserRecord(&customer.User)
	if err != nil {
		return nil, err
	}

	return &pkg.DefaultResponse{
		Success: "success",
		Message: "user details updated successfully",
	}, nil
}

func (u userAppHandler) GetAuthenticatedUser(ctx context.Context) (*models.Customer, error) {
	claims, err := u.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	customer, err := u.allRepository.UserRepository.GetCustomerByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	return customer, nil
}
