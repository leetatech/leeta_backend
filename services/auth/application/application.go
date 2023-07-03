package application

import (
	"github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/models"
	"go.uber.org/zap"
	"time"
)

type authAppHandler struct {
	tokenHandler  library.TokenHandler
	encryptor     library.EncryptorManager
	idGenerator   library.IDGenerator
	otpGenerator  library.OtpGenerator
	logger        *zap.Logger
	allRepository library.Repositories
}

type AuthApplication interface {
	SignUp(request domain.SignUpRequest) (*domain.DefaultSigningResponse, error)
	CreateOTP(request domain.OTPRequest) (*library.DefaultResponse, error)
	EarlyAccess(request models.EarlyAccess) (*library.DefaultResponse, error)
}

func NewAuthApplication(tokenHandler library.TokenHandler, logger *zap.Logger, allRepository library.Repositories) AuthApplication {
	return &authAppHandler{
		tokenHandler:  tokenHandler,
		encryptor:     library.NewEncryptor(),
		idGenerator:   library.NewIDGenerator(),
		otpGenerator:  library.NewOTPGenerator(),
		logger:        logger,
		allRepository: allRepository,
	}
}

func (a authAppHandler) SignUp(request domain.SignUpRequest) (*domain.DefaultSigningResponse, error) {
	hashedPassword, err := a.passwordValidationEncryption(request.Password)
	if err != nil {
		a.logger.Error("Password Validation", zap.Error(err))
		return nil, err
	}

	request.Password = hashedPassword
	if models.IsValidUserCategory(request.UserType) {
		switch request.UserType {
		case models.VendorCategory:

			return a.vendorSignUP(request)
		}
	}

	// TODO: send email to user
	return nil, nil
}

func (a authAppHandler) CreateOTP(request domain.OTPRequest) (*library.DefaultResponse, error) {
	expirationDuration := time.Duration(5) * time.Minute

	otpResponse := models.Verification{
		ID:        a.idGenerator.Generate(),
		Code:      a.otpGenerator.Generate(),
		Topic:     request.Topic,
		Type:      request.Type,
		Target:    request.Target,
		ExpiresAt: time.Now().Add(expirationDuration).Unix(),
		Timestamp: time.Now().Unix(),
	}

	err := a.allRepository.AuthRepository.CreateOTP(otpResponse)
	if err != nil {
		return nil, err
	}

	return &library.DefaultResponse{Success: "success", Message: "OTP created"}, nil
}

func (a authAppHandler) EarlyAccess(request models.EarlyAccess) (*library.DefaultResponse, error) {
	err := a.allRepository.AuthRepository.EarlyAccess(request)
	if err != nil {
		a.logger.Error("EarlyAccess", zap.Any(leetError.ErrorType(leetError.DatabaseError), err), zap.Any(leetError.ErrorType(leetError.DatabaseError), leetError.ErrorMessage(leetError.DatabaseError)))
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	// TODO: send email to user

	return &library.DefaultResponse{Success: "success", Message: "Early Access created"}, nil
}
