package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/encrypto"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/pkg/idgenerator"
	"github.com/leetatech/leeta_backend/pkg/jwtmiddleware"
	"github.com/leetatech/leeta_backend/pkg/mailer"
	"github.com/leetatech/leeta_backend/pkg/otp"
	"github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/auth/infrastructure"
	"github.com/leetatech/leeta_backend/services/models"
	"time"
)

type authAppHandler struct {
	jwtManager        jwtmiddleware.Manager
	encryptor         encrypto.Manager
	idGenerator       idgenerator.Generator
	otpGenerator      otp.Generator
	mailer            mailer.Client
	domain            string
	repositoryManager pkg.RepositoryManager
}

type Auth interface {
	SignUp(ctx context.Context, request domain.SignupRequest) (*domain.DefaultSigningResponse, error)
	RequestOTP(ctx context.Context, request domain.EmailRequestBody) (*pkg.DefaultResponse, error)
	EarlyAccess(ctx context.Context, request models.EarlyAccess) (*pkg.DefaultResponse, error)
	SignIn(ctx context.Context, request domain.SigningRequest) (*domain.DefaultSigningResponse, error)
	ForgotPassword(ctx context.Context, request domain.EmailRequestBody) (*pkg.DefaultResponse, error)
	ValidateOTP(ctx context.Context, request domain.OTPValidationRequest) (*pkg.DefaultResponse, error)
	CreateNewPassword(ctx context.Context, request domain.CreateNewPasswordRequest) (*domain.DefaultSigningResponse, error)
	AdminSignUp(ctx context.Context, request domain.AdminSignUpRequest) (*domain.DefaultSigningResponse, error)
	ReceiveGuestToken(request domain.ReceiveGuestRequest) (*domain.ReceiveGuestResponse, error)
	UpdateGuestRecord(ctx context.Context, request models.Guest) (*pkg.DefaultResponse, error)
	GetGuestRecord(ctx context.Context, deviceId string) (models.Guest, error)
}

func New(request pkg.ApplicationContext) Auth {
	return &authAppHandler{
		jwtManager:        request.JwtManager,
		encryptor:         encrypto.New(),
		idGenerator:       idgenerator.New(),
		otpGenerator:      otp.New(),
		mailer:            request.Mailer,
		domain:            request.Domain,
		repositoryManager: request.RepositoryManager,
	}
}

func (a authAppHandler) SignUp(ctx context.Context, request domain.SignupRequest) (*domain.DefaultSigningResponse, error) {
	hashedPassword, err := a.validateAndEncryptPassword(request.Password)
	if err != nil {
		return nil, fmt.Errorf("error validating and encrypting password on signup: %w", err)
	}
	request.Password = hashedPassword

	// trim email spaces
	trimmedEmail := request.TrimEmailSpace()
	request.Email = trimmedEmail

	category, err := models.SetUserCategory(request.UserType)
	if err != nil {
		return nil, err
	}

	switch category {
	case models.VendorCategory:
		return a.vendorSignUP(ctx, request)

	case models.CustomerCategory:
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

	err := a.repositoryManager.AuthRepository.CreateOTP(ctx, otpResponse)
	if err != nil {
		return nil, err
	}

	return &pkg.DefaultResponse{Success: "success", Message: otpResponse.Code}, nil
}

func (a authAppHandler) EarlyAccess(ctx context.Context, request models.EarlyAccess) (*pkg.DefaultResponse, error) {
	request.Timestamp = time.Now().Unix()
	err := a.repositoryManager.AuthRepository.SaveEarlyAccess(ctx, request)
	if err != nil {
		return nil, errs.Body(errs.DatabaseError, fmt.Errorf("error saving early access: %w", err))
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
	// trim email spaces
	trimmedEmail := request.TrimEmailSpace()
	request.Email = trimmedEmail

	category, err := models.SetUserCategory(request.UserType)
	if err != nil {
		return nil, err
	}
	switch category {
	case models.VendorCategory:
		return a.vendorSignIN(ctx, request)
	case models.AdminCategory:
		return a.adminSignIN(ctx, request)
	case models.CustomerCategory:
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
		return nil, errs.Body(errs.ForgotPasswordError, err)
	}
	return &pkg.DefaultResponse{Success: "success", Message: "An email with an OTP has been sent to you"}, nil
}

func (a authAppHandler) sendOTP(ctx context.Context, request domain.EmailRequestBody) error {
	// get user by email
	user, err := a.repositoryManager.AuthRepository.UserByEmail(ctx, request.Email)
	if err != nil {
		return errs.Body(errs.UserNotFoundError, fmt.Errorf("error getting user by email: %w", err))
	}

	// check if user otp exists
	verification, err := a.repositoryManager.AuthRepository.FindUnvalidatedVerificationByTarget(ctx, request.Email)
	if err != nil && !errors.Is(err, infrastructure.ErrItemNotFound) {
		return errs.Body(errs.DatabaseError, err)
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
	verification, err := a.repositoryManager.AuthRepository.FindUnvalidatedVerificationByTarget(ctx, request.Target)
	if err != nil {
		return nil, fmt.Errorf("error getting unvalidated verification by target when validating otp: %w", err)
	}
	if verification.Validated {
		return nil, errs.Body(errs.TokenValidationError, fmt.Errorf("otp has already been validated"))
	}

	if verification.Code != request.Code {
		return nil, errs.Body(errs.TokenValidationError, errors.New("invalid otp"))
	}

	if time.Unix(verification.ExpiresAt, 0).Before(time.Now()) {
		return nil, errs.Body(errs.TokenValidationError, errors.New("expired otp"))
	}

	err = a.repositoryManager.AuthRepository.ValidateOTP(ctx, verification.ID)
	if err != nil {
		return nil, fmt.Errorf("error validating otp: %w", err)
	}

	err = a.repositoryManager.AuthRepository.SetEmailVerificationStatus(ctx, verification.Target, true)
	if err != nil {
		return nil, fmt.Errorf("error setting email verification status")
	}

	return &pkg.DefaultResponse{Success: "success", Message: "OTP validated"}, nil
}

func (a authAppHandler) CreateNewPassword(ctx context.Context, request domain.CreateNewPasswordRequest) (*domain.DefaultSigningResponse, error) {

	if request.Password != request.ConfirmPassword {
		return nil, errs.Body(errs.PasswordValidationError, errors.New("password and confirm password don't match"))
	}

	vendor, err := a.repositoryManager.AuthRepository.UserByEmail(ctx, request.Email)
	if err != nil {
		return nil, errs.Body(errs.UserNotFoundError, fmt.Errorf("error getting user by email: %w", err))
	}

	return a.createNewPassword(ctx, vendor.ID, request.Password)

}

func (a authAppHandler) AdminSignUp(ctx context.Context, request domain.AdminSignUpRequest) (*domain.DefaultSigningResponse, error) {
	err := a.encryptor.ValidateEmailFormat(request.Email)
	if err != nil {
		return nil, fmt.Errorf("error validating email format: %w", err)
	}

	err = a.encryptor.ValidateDomain(request.Email, a.domain)
	if err != nil {
		return nil, fmt.Errorf("error validating domain on admin sign up: %w", err)
	}

	hashedPassword, err := a.validateAndEncryptPassword(request.Password)
	if err != nil {
		return nil, fmt.Errorf("error validating and encrypting password on admin sign up: %w", err)
	}
	request.Password = hashedPassword

	return a.adminSignUp(ctx, request)
}

func (a authAppHandler) ReceiveGuestToken(request domain.ReceiveGuestRequest) (*domain.ReceiveGuestResponse, error) {
	ctx := context.Background()

	// check if guest device id already exist. if it does then there is already an assigned guest id
	guestRecord, err := a.repositoryManager.AuthRepository.GuestRecord(ctx, request.DeviceID)
	if err != nil {
		if !errors.Is(err, infrastructure.ErrItemNotFound) {
			return nil, errs.Body(errs.InternalError, fmt.Errorf("error when searching for guest record %w", err))
		}
	}

	if guestRecord.ID == "" {
		guestID := a.idGenerator.Generate()
		// store guest information
		guestRecord = models.Guest{
			ID:       guestID,
			DeviceID: request.DeviceID,
			Location: request.Location,
		}

		if err := a.repositoryManager.AuthRepository.CreateGuestRecord(context.Background(), guestRecord); err != nil {
			return nil, errs.Body(errs.InternalError, fmt.Errorf("error creating guest record %w", err))
		}
	}

	tokenString, err := a.jwtManager.BuildAuthResponse("", guestRecord.ID, request.DeviceID, models.GuestCategory)
	if err != nil {
		return nil, errs.Body(errs.InternalError, fmt.Errorf("error building token response %w", err))
	}

	return &domain.ReceiveGuestResponse{
		SessionID: guestRecord.ID,
		DeviceID:  request.DeviceID,
		Token:     tokenString,
	}, nil
}

func (a authAppHandler) UpdateGuestRecord(ctx context.Context, request models.Guest) (*pkg.DefaultResponse, error) {
	guestRecord, err := a.repositoryManager.AuthRepository.GuestRecord(ctx, request.DeviceID)
	if err != nil {
		if !errors.Is(err, infrastructure.ErrItemNotFound) {
			return nil, errs.Body(errs.InternalError, fmt.Errorf("error when searching for guest record %w", err))
		}
	}

	guestRecord.FirstName = request.FirstName
	guestRecord.LastName = request.LastName
	guestRecord.Number = request.Number
	guestRecord.Email = request.Email
	guestRecord.Address.State = request.Address.State
	guestRecord.Address.City = request.Address.City
	guestRecord.Address.LGA = request.Address.LGA
	guestRecord.Address.FullAddress = request.Address.FullAddress
	guestRecord.Address.ClosestLandmark = request.Address.ClosestLandmark
	guestRecord.Address.Coordinates = request.Address.Coordinates
	guestRecord.Address.Verified = request.Address.Verified

	err = a.repositoryManager.AuthRepository.UpdateGuestRecord(ctx, guestRecord)
	if err != nil {
		return nil, err
	}
	return &pkg.DefaultResponse{Success: "success", Message: "Guest record updated"}, nil
}

func (a authAppHandler) GetGuestRecord(ctx context.Context, deviceId string) (models.Guest, error) {
	return a.repositoryManager.AuthRepository.GuestRecord(ctx, deviceId)
}
