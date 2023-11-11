package application

import (
	"context"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/mailer"
	"github.com/leetatech/leeta_backend/services/library/models"
	"github.com/leetatech/leeta_backend/services/user/domain"
	"go.uber.org/zap"
	"time"
)

type userAppHandler struct {
	tokenHandler  library.TokenHandler
	encryptor     library.EncryptorManager
	idGenerator   library.IDGenerator
	otpGenerator  library.OtpGenerator
	logger        *zap.Logger
	EmailClient   mailer.MailerClient
	allRepository library.Repositories
}

type UserApplication interface {
	VendorVerification(ctx context.Context, request domain.VendorVerificationRequest) (*library.DefaultResponse, error)
	AddVendorByAdmin(ctx context.Context, request domain.VendorVerificationRequest) (*library.DefaultResponse, error)
}

func NewUserApplication(request library.DefaultApplicationRequest) UserApplication {
	return &userAppHandler{
		tokenHandler:  request.TokenHandler,
		encryptor:     library.NewEncryptor(),
		idGenerator:   library.NewIDGenerator(),
		otpGenerator:  library.NewOTPGenerator(),
		logger:        request.Logger,
		EmailClient:   request.EmailClient,
		allRepository: request.AllRepository,
	}
}

func (u userAppHandler) VendorVerification(ctx context.Context, request domain.VendorVerificationRequest) (*library.DefaultResponse, error) {
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

	return &library.DefaultResponse{Success: "success", Message: "Business successfully registered"}, nil
}

func (u userAppHandler) AddVendorByAdmin(ctx context.Context, request domain.VendorVerificationRequest) (*library.DefaultResponse, error) {
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
		Identity: request.Identity,
		AdminID:  claims.UserID,
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

	return &library.DefaultResponse{Success: "success", Message: "Business successfully registered"}, nil
}
