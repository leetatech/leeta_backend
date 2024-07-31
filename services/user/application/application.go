package application

import (
	"context"
	"time"

	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/encrypto"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/pkg/idgenerator"
	"github.com/leetatech/leeta_backend/pkg/jwtmiddleware"
	"github.com/leetatech/leeta_backend/pkg/mailer/aws"
	"github.com/leetatech/leeta_backend/pkg/otp"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/user/domain"
)

type userAppHandler struct {
	jwtManager    jwtmiddleware.Manager
	encryptor     encrypto.Manager
	idGenerator   idgenerator.Generator
	otpGenerator  otp.Generator
	EmailClient   aws.MailClient
	allRepository pkg.RepositoryManager
}

type UserApplication interface {
	VendorVerification(ctx context.Context, request domain.VendorVerificationRequest) (*pkg.DefaultResponse, error)
	AddVendorByAdmin(ctx context.Context, request domain.VendorVerificationRequest) (*pkg.DefaultResponse, error)
	Data(ctx context.Context) (*models.Customer, error)
	UpdateRecord(ctx context.Context, request models.User) (*pkg.DefaultResponse, error)
}

func New(request pkg.ApplicationContext) UserApplication {
	return &userAppHandler{
		jwtManager:    request.JwtManager,
		encryptor:     encrypto.New(),
		idGenerator:   idgenerator.New(),
		otpGenerator:  otp.New(),
		EmailClient:   request.MailClient,
		allRepository: request.RepositoryManager,
	}
}

func (u *userAppHandler) VendorVerification(ctx context.Context, request domain.VendorVerificationRequest) (*pkg.DefaultResponse, error) {
	claims, err := u.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return nil, errs.Body(errs.ErrorUnauthorized, err)
	}

	if claims.Role != models.VendorCategory {
		return nil, errs.Body(errs.ErrorUnauthorized, err)
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

func (u *userAppHandler) AddVendorByAdmin(ctx context.Context, request domain.VendorVerificationRequest) (*pkg.DefaultResponse, error) {
	claims, err := u.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return nil, errs.Body(errs.ErrorUnauthorized, err)
	}
	if claims.Role != models.AdminCategory {
		return nil, errs.Body(errs.ErrorUnauthorized, err)
	}
	_, err = u.allRepository.AuthRepository.AdminByEmail(ctx, claims.Email)
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

func (u *userAppHandler) UpdateRecord(ctx context.Context, request models.User) (*pkg.DefaultResponse, error) {
	claims, err := u.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return nil, errs.Body(errs.ErrorUnauthorized, err)
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

func (u *userAppHandler) Data(ctx context.Context) (*models.Customer, error) {
	claims, err := u.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return nil, errs.Body(errs.ErrorUnauthorized, err)
	}

	customer, err := u.allRepository.UserRepository.GetCustomerByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	return customer, nil
}
