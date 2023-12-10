package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/mailer"
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
	EmailClient   mailer.MailerClient
	Domain        string
	allRepository library.Repositories
}

type AuthApplication interface {
	SignUp(ctx context.Context, request domain.SignupRequest) (*domain.DefaultSigningResponse, error)
	CreateOTP(ctx context.Context, request domain.OTPRequest) (*library.DefaultResponse, error)
	EarlyAccess(ctx context.Context, request models.EarlyAccess) (*library.DefaultResponse, error)
	SignIn(ctx context.Context, request domain.SigningRequest) (*domain.DefaultSigningResponse, error)
	ForgotPassword(ctx context.Context, request domain.ForgotPasswordRequest) (*library.DefaultResponse, error)
	ValidateOTP(ctx context.Context, request domain.OTPValidationRequest) (*library.DefaultResponse, error)
	ResetPassword(ctx context.Context, request domain.ResetPasswordRequest) (*domain.DefaultSigningResponse, error)
	AdminSignUp(ctx context.Context, request domain.AdminSignUpRequest) (*domain.DefaultSigningResponse, error)
}

func NewAuthApplication(request library.DefaultApplicationRequest) AuthApplication {
	return &authAppHandler{
		tokenHandler:  request.TokenHandler,
		encryptor:     library.NewEncryptor(),
		idGenerator:   library.NewIDGenerator(),
		otpGenerator:  library.NewOTPGenerator(),
		logger:        request.Logger,
		EmailClient:   request.EmailClient,
		Domain:        request.Domain,
		allRepository: request.AllRepository,
	}
}

func (a authAppHandler) SignUp(ctx context.Context, request domain.SignupRequest) (*domain.DefaultSigningResponse, error) {
	hashedPassword, err := a.passwordValidationEncryption(request.Password)
	if err != nil {
		a.logger.Error("Password Validation", zap.Error(err))
		return nil, err
	}
	request.Password = hashedPassword

	category, err := models.SetUserCategory(request.UserType)
	if err != nil {
		return nil, err
	}

	switch category {
	case models.VendorCategory:
		return a.vendorSignUP(ctx, request)

	case models.BuyerCategory:
		return a.customerSignUP(ctx, request)
	}

	return nil, nil
}

func (a authAppHandler) CreateOTP(ctx context.Context, request domain.OTPRequest) (*library.DefaultResponse, error) {
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

	err := a.allRepository.AuthRepository.CreateOTP(ctx, otpResponse)
	if err != nil {
		return nil, err
	}

	return &library.DefaultResponse{Success: "success", Message: otpResponse.Code}, nil
}

func (a authAppHandler) EarlyAccess(ctx context.Context, request models.EarlyAccess) (*library.DefaultResponse, error) {
	request.Timestamp = time.Now().Unix()
	err := a.allRepository.AuthRepository.EarlyAccess(ctx, request)
	if err != nil {
		a.logger.Error("EarlyAccess", zap.Any(leetError.ErrorType(leetError.DatabaseError), err), zap.Any(leetError.ErrorType(leetError.DatabaseError), leetError.ErrorMessage(leetError.DatabaseError)))
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	message := models.Message{
		ID:         a.idGenerator.Generate(),
		Target:     request.Email,
		TemplateID: library.EarlyAccessEmailTemplateID,
		DataMap: map[string]string{
			"URL": "https://deploy-preview-3--gleeful-palmier-8efb17.netlify.app/",
		},
		Ts: time.Now().Unix(),
	}
	err = a.sendEmail(message)
	if err != nil {
		return nil, err
	}

	return &library.DefaultResponse{Success: "success", Message: "Early Access created"}, nil
}

func (a authAppHandler) SignIn(ctx context.Context, request domain.SigningRequest) (*domain.DefaultSigningResponse, error) {
	category, err := models.SetUserCategory(request.UserType)
	if err != nil {
		return nil, err
	}
	switch category {
	case models.VendorCategory:
		return a.vendorSignIN(ctx, request)
	case models.AdminCategory:
		return a.adminSignIN(ctx, request)
	case models.BuyerCategory:
		return a.customerSignIN(ctx, request)
	}

	return nil, nil
}

func (a authAppHandler) ForgotPassword(ctx context.Context, request domain.ForgotPasswordRequest) (*library.DefaultResponse, error) {
	// get customer by email
	customer, err := a.allRepository.AuthRepository.GetCustomerByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	// get customer category from roles
	identity, err := a.allRepository.AuthRepository.GetIdentityByCustomerID(ctx, customer.ID)
	if err != nil {
		return nil, err
	}

	category := identity.Role

	var firstName, lastName string
	switch category {
	case models.VendorCategory:
		vendor, err := a.allRepository.AuthRepository.GetVendorByEmail(ctx, request.Email)
		if err != nil {
			a.logger.Error("SignIn", zap.Any(leetError.ErrorType(leetError.UserNotFoundError), err), zap.Any("email", request.Email))
			return nil, leetError.ErrorResponseBody(leetError.UserNotFoundError, err)
		}
		firstName = vendor.FirstName
		lastName = vendor.LastName
	}

	requestOTP := domain.OTPRequest{
		Topic:        "ForgotPassword",
		Type:         models.EMAIL,
		Target:       request.Email,
		UserCategory: category,
	}

	otpResponse, err := a.CreateOTP(ctx, requestOTP)
	if err != nil {
		return nil, err
	}

	message := models.Message{
		ID:         a.idGenerator.Generate(),
		Target:     request.Email,
		TemplateID: library.ForgotPasswordEmailTemplateID,
		DataMap: map[string]string{
			"FirstName": firstName,
			"LastName":  lastName,
			"OTP":       otpResponse.Message,
		},
		Ts: time.Now().Unix(),
	}
	err = a.sendEmail(message)
	if err != nil {
		return nil, err
	}

	return &library.DefaultResponse{Success: "success", Message: "OTP created"}, nil
}

func (a authAppHandler) ValidateOTP(ctx context.Context, request domain.OTPValidationRequest) (*library.DefaultResponse, error) {
	verification, err := a.allRepository.AuthRepository.GetOTPForValidation(ctx, request.Target)
	if err != nil {
		a.logger.Error("ValidateOTP", zap.String("target", request.Target), zap.Error(err))
		return nil, err
	}
	if verification.Validated {
		newErr := errors.New("already validated otp")
		a.logger.Error("ValidateOTP", zap.String(leetError.ErrorType(leetError.TokenValidationError), fmt.Sprintf("%s: %s", "target", request.Target)), zap.Error(newErr))

		return nil, leetError.ErrorResponseBody(leetError.TokenValidationError, leetError.ErrorResponseBody(leetError.TokenValidationError, newErr))
	}

	if verification.Code != request.Code {
		newErr := errors.New("invalid otp")
		a.logger.Error("ValidateOTP", zap.String(leetError.ErrorType(leetError.TokenValidationError), fmt.Sprintf("%s: %s", "target", request.Target)), zap.Error(leetError.ErrorResponseBody(leetError.TokenValidationError, newErr)))

		return nil, leetError.ErrorResponseBody(leetError.TokenValidationError, newErr)
	}

	if time.Unix(verification.ExpiresAt, 0).Before(time.Now()) {
		newErr := errors.New("expired otp")
		a.logger.Error("ValidateOTP", zap.String(leetError.ErrorType(leetError.TokenValidationError), fmt.Sprintf("%s: %s", "target", request.Target)), zap.Error(leetError.ErrorResponseBody(leetError.TokenValidationError, leetError.ErrorResponseBody(leetError.TokenValidationError, newErr))))

		return nil, leetError.ErrorResponseBody(leetError.TokenValidationError, newErr)
	}

	err = a.allRepository.AuthRepository.ValidateOTP(ctx, verification.ID)
	if err != nil {
		a.logger.Error("store validating verification", zap.Error(err), zap.String("verification_id", verification.ID))
		return nil, err
	}

	return &library.DefaultResponse{Success: "success", Message: "OTP validated"}, nil
}

func (a authAppHandler) ResetPassword(ctx context.Context, request domain.ResetPasswordRequest) (*domain.DefaultSigningResponse, error) {

	if request.Password != request.ConfirmPassword {
		a.logger.Error("ResetPassword", zap.Any(leetError.ErrorType(leetError.PasswordValidationError), errors.New("password and confirm password don't match")))

		return nil, leetError.ErrorResponseBody(leetError.PasswordValidationError, errors.New("password and confirm password don't match"))
	}

	category, err := models.SetUserCategory(request.UserCategory)
	if err != nil {
		return nil, err
	}
	switch category {
	case models.VendorCategory:
		vendor, err := a.allRepository.AuthRepository.GetVendorByEmail(ctx, request.Email)
		if err != nil {
			a.logger.Error("ResetPassword", zap.Any(leetError.ErrorType(leetError.UserNotFoundError), err), zap.Any("email", request.Email))
			return nil, leetError.ErrorResponseBody(leetError.UserNotFoundError, err)
		}

		return a.resetPassword(ctx, vendor.ID, vendor.Email.Address, request.Password, request.UserCategory)
	}

	return nil, nil
}

func (a authAppHandler) AdminSignUp(ctx context.Context, request domain.AdminSignUpRequest) (*domain.DefaultSigningResponse, error) {
	err := a.encryptor.IsValidEmailFormat(request.Email)
	if err != nil {
		a.logger.Error("AdminSignUp", zap.Error(err))
		return nil, err
	}

	err = a.encryptor.IsLeetaDomain(request.Email, a.Domain)
	if err != nil {
		a.logger.Error("AdminSignUp", zap.Error(err))
		return nil, err
	}

	hashedPassword, err := a.passwordValidationEncryption(request.Password)
	if err != nil {
		a.logger.Error("Password Validation", zap.Error(err))
		return nil, err
	}
	request.Password = hashedPassword

	return a.adminSignUp(ctx, request)
}
