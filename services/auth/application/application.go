package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/mailer"
	"github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/auth/infrastructure"
	"github.com/leetatech/leeta_backend/services/models"
	"go.uber.org/zap"
	"time"
)

type authAppHandler struct {
	tokenHandler  pkg.TokenHandler
	encryptor     pkg.EncryptorManager
	idGenerator   pkg.IDGenerator
	otpGenerator  pkg.OtpGenerator
	logger        *zap.Logger
	EmailClient   mailer.MailerClient
	Domain        string
	allRepository pkg.Repositories
}

type AuthApplication interface {
	SignUp(ctx context.Context, request domain.SignupRequest) (*domain.DefaultSigningResponse, error)
	RequestOTP(ctx context.Context, request domain.EmailRequestBody) (*pkg.DefaultResponse, error)
	EarlyAccess(ctx context.Context, request models.EarlyAccess) (*pkg.DefaultResponse, error)
	SignIn(ctx context.Context, request domain.SigningRequest) (*domain.DefaultSigningResponse, error)
	ForgotPassword(ctx context.Context, request domain.EmailRequestBody) (*pkg.DefaultResponse, error)
	ValidateOTP(ctx context.Context, request domain.OTPValidationRequest) (*pkg.DefaultResponse, error)
	CreateNewPassword(ctx context.Context, request domain.CreateNewPasswordRequest) (*domain.DefaultSigningResponse, error)
	AdminSignUp(ctx context.Context, request domain.AdminSignUpRequest) (*domain.DefaultSigningResponse, error)
	ReceiveGuestToken(request domain.ReceiveGuestRequest) (*domain.ReceiveGuestResponse, error)
}

func NewAuthApplication(request pkg.DefaultApplicationRequest) AuthApplication {
	return &authAppHandler{
		tokenHandler:  request.TokenHandler,
		encryptor:     pkg.NewEncryptor(),
		idGenerator:   pkg.NewIDGenerator(),
		otpGenerator:  pkg.NewOTPGenerator(),
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

func (a authAppHandler) createOTP(ctx context.Context, request domain.OTPRequest) (*pkg.DefaultResponse, error) {
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

	return &pkg.DefaultResponse{Success: "success", Message: otpResponse.Code}, nil
}

func (a authAppHandler) EarlyAccess(ctx context.Context, request models.EarlyAccess) (*pkg.DefaultResponse, error) {
	request.Timestamp = time.Now().Unix()
	err := a.allRepository.AuthRepository.EarlyAccess(ctx, request)
	if err != nil {
		a.logger.Error("EarlyAccess", zap.Any(leetError.ErrorType(leetError.DatabaseError), err), zap.Any(leetError.ErrorType(leetError.DatabaseError), leetError.ErrorMessage(leetError.DatabaseError)))
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	message := models.Message{
		ID:         a.idGenerator.Generate(),
		Target:     request.Email,
		TemplateID: pkg.EarlyAccessEmailTemplateID,
		DataMap: map[string]string{
			"URL": "https://deploy-preview-3--gleeful-palmier-8efb17.netlify.app/",
		},
		Ts: time.Now().Unix(),
	}
	err = a.sendEmail(message)
	if err != nil {
		return nil, err
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Early Access created"}, nil
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

func (a authAppHandler) ForgotPassword(ctx context.Context, request domain.EmailRequestBody) (*pkg.DefaultResponse, error) {
	if err := a.sendOTP(ctx, request); err != nil {
		return nil, err
	}
	return &pkg.DefaultResponse{Success: "success", Message: "An email with OTP to reset your password has been sent to you"}, nil
}

func (a authAppHandler) RequestOTP(ctx context.Context, request domain.EmailRequestBody) (*pkg.DefaultResponse, error) {
	if err := a.sendOTP(ctx, request); err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ForgotPasswordError, err)
	}
	return &pkg.DefaultResponse{Success: "success", Message: "An email with an OTP has been sent to you"}, nil
}

func (a authAppHandler) sendOTP(ctx context.Context, request domain.EmailRequestBody) error {
	// get user by email
	user, err := a.allRepository.AuthRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		a.logger.Error("error getting user by email", zap.Any(leetError.ErrorType(leetError.UserNotFoundError), err), zap.Any("email", request.Email))
		return leetError.ErrorResponseBody(leetError.UserNotFoundError, err)
	}

	// check if user otp exists
	verification, err := a.allRepository.AuthRepository.GetOTPForValidation(ctx, request.Email)
	if err != nil && !errors.Is(err, infrastructure.ErrItemNotFound) {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	requestOTP := domain.OTPRequest{
		Topic:  "ForgotPassword",
		Type:   models.EMAIL,
		Target: request.Email,
	}
	var OTP string
	if isVerificationValid := verification.VerifyCodeValidity(); !isVerificationValid {
		response, err := a.createOTP(ctx, requestOTP)
		if err != nil {
			return err
		}
		OTP = response.Message
	} else {
		OTP = verification.Code
	}

	message := models.Message{
		ID:         a.idGenerator.Generate(),
		Target:     request.Email,
		TemplateID: pkg.ForgotPasswordEmailTemplateID,
		DataMap: map[string]string{
			"FirstName": user.FirstName,
			"LastName":  user.LastName,
			"OTP":       OTP,
		},
		Ts: time.Now().Unix(),
	}
	err = a.sendEmail(message)
	if err != nil {
		return err
	}

	return nil
}

func (a authAppHandler) ValidateOTP(ctx context.Context, request domain.OTPValidationRequest) (*pkg.DefaultResponse, error) {
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

	err = a.allRepository.AuthRepository.UpdateEmailVerify(ctx, verification.Target, true)
	if err != nil {
		a.logger.Error("error validating user email", zap.Error(err), zap.String("verification_email", verification.Target))
		return nil, err
	}

	return &pkg.DefaultResponse{Success: "success", Message: "OTP validated"}, nil
}

func (a authAppHandler) CreateNewPassword(ctx context.Context, request domain.CreateNewPasswordRequest) (*domain.DefaultSigningResponse, error) {

	if request.Password != request.ConfirmPassword {
		a.logger.Error("CreateNewPassword", zap.Any(leetError.ErrorType(leetError.PasswordValidationError), errors.New("password and confirm password don't match")))

		return nil, leetError.ErrorResponseBody(leetError.PasswordValidationError, errors.New("password and confirm password don't match"))
	}

	vendor, err := a.allRepository.AuthRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		a.logger.Error("CreateNewPassword", zap.Any(leetError.ErrorType(leetError.UserNotFoundError), err), zap.Any("email", request.Email))
		return nil, leetError.ErrorResponseBody(leetError.UserNotFoundError, err)
	}

	return a.createNewPassword(ctx, vendor.ID, vendor.Email.Address, request.Password)

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

func (a authAppHandler) ReceiveGuestToken(request domain.ReceiveGuestRequest) (*domain.ReceiveGuestResponse, error) {
	sessionID := a.idGenerator.Generate()
	tokenString, err := a.tokenHandler.BuildAuthResponse("", request.DeviceID, sessionID, models.GuestCatergory)
	if err != nil {
		return nil, err
	}

	return &domain.ReceiveGuestResponse{
		SessionID: sessionID,
		DeviceID:  request.DeviceID,
		Token:     tokenString,
	}, nil
}
