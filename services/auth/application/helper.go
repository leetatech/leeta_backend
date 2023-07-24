package application

import (
	"errors"
	"github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"sync"
	"time"
)

func (a authAppHandler) passwordValidationEncryption(password string) (string, error) {
	err := a.encryptor.IsValidPassword(password)
	if err != nil {
		a.logger.Error("passwordValidationEncryption", zap.Error(err))
		return "", err
	}
	passByte, err := a.encryptor.GenerateFromPasscode(password)
	if err != nil {
		return "", leetError.ErrorResponseBody(leetError.EncryptionError, err)
	}

	return string(passByte), nil
}

func (a authAppHandler) vendorSignUP(request domain.SigningRequest) (*domain.DefaultSigningResponse, error) {
	_, err := a.allRepository.AuthRepository.GetVendorByEmail(request.Email)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			timestamp := time.Now().Unix()

			vendor := models.Vendor{
				ID: a.idGenerator.Generate(),
				Email: models.Email{
					Address: request.Email,
				},
				Timestamp: timestamp,
			}
			err = a.allRepository.AuthRepository.CreateVendor(vendor)
			if err != nil {
				return nil, err
			}

			identity := models.Identity{
				ID:         a.idGenerator.Generate(),
				CustomerID: vendor.ID,
				Role:       models.VendorCategory,
				Credentials: []models.Credentials{
					{
						Type:            models.CredentialsTypeLogin,
						Password:        request.Password,
						Status:          models.CredentialStatusActive,
						StatusTimestamp: timestamp,
						Timestamp:       timestamp,
					},
				},
			}
			err = a.allRepository.AuthRepository.CreateIdentity(identity)
			if err != nil {
				return nil, err
			}

			response, err := a.tokenHandler.BuildAuthResponse(request.Email, vendor.ID, a.idGenerator.Generate(), request.UserType)
			if err != nil {
				a.logger.Error("SignUp", zap.Any("BuildAuthResponse", leetError.ErrorResponseBody(leetError.TokenGenerationError, err)))
				return nil, leetError.ErrorResponseBody(leetError.TokenGenerationError, err)
			}

			err = a.accountVerification(vendor.ID, vendor.Email.Address, vendor.FirstName, vendor.LastName)
			if err != nil {
				return nil, err
			}

			return &domain.DefaultSigningResponse{AuthToken: response}, nil

		default:
			return nil, err
		}
	}

	a.logger.Error("vendorSignUP", zap.Error(leetError.ErrorResponseBody(leetError.DuplicateUserError, nil)))
	return nil, leetError.ErrorResponseBody(leetError.DuplicateUserError, nil)
}

func (a authAppHandler) accountVerification(userID, target, firstName, lastName string) error {
	requestOTP := domain.OTPRequest{
		Topic:        "Sign Up",
		Type:         models.EMAIL,
		Target:       target,
		UserCategory: models.VendorCategory,
	}
	otpResponse, err := a.CreateOTP(requestOTP)
	if err != nil {
		a.logger.Error("SignUp", zap.Any("CreateOTP", err))
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	message := models.Message{
		ID:         a.idGenerator.Generate(),
		UserID:     userID,
		Target:     target,
		TemplateID: library.SignUpEmailTemplateID,
		DataMap: map[string]string{
			"OTP": otpResponse.Message,
		},
		Ts: time.Now().Unix(),
	}
	err = a.sendEmail(message)
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}

func (a authAppHandler) vendorSignIN(request domain.SigningRequest) (*domain.DefaultSigningResponse, error) {
	vendor, err := a.allRepository.AuthRepository.GetVendorByEmail(request.Email)
	if err != nil {
		a.logger.Error("SignIn", zap.Any(leetError.ErrorType(leetError.UserNotFoundError), err), zap.Any("email", request.Email))
		return nil, leetError.ErrorResponseBody(leetError.UserNotFoundError, err)
	}

	identity, err := a.allRepository.AuthRepository.GetIdentityByCustomerID(vendor.ID)
	if err != nil {
		a.logger.Error("SignIn", zap.Any(leetError.ErrorType(leetError.IdentityNotFoundError), err), zap.Any("vendor_id", vendor.ID))
		return nil, leetError.ErrorResponseBody(leetError.IdentityNotFoundError, err)
	}

	switch vendor.Status {
	case models.Locked, models.Exited, models.Rejected:
		a.logger.Error("SignIn", zap.Any(leetError.ErrorType(leetError.UserLockedError), err), zap.Any(leetError.ErrorType(leetError.UserLockedError), leetError.ErrorMessage(leetError.UserLockedError)))
		return nil, leetError.ErrorResponseBody(leetError.UserLockedError, err)
	}

	err = a.processLoginPasswordValidation(request, identity)
	if err != nil {
		return nil, err
	}

	response, err := a.tokenHandler.BuildAuthResponse(request.Email, vendor.ID, identity.ID, request.UserType)
	if err != nil {
		a.logger.Error("SignIn", zap.Any("BuildAuthResponse", leetError.ErrorResponseBody(leetError.TokenGenerationError, err)))
		return nil, leetError.ErrorResponseBody(leetError.TokenGenerationError, err)
	}
	//message := models.Message{
	//	ID:         a.idGenerator.Generate(),
	//	UserID:     vendor.ID,
	//	Target:     vendor.Email.Address,
	//	TemplateID: library.EarlyAccessEmailTemplateID,
	//	Ts:         time.Now().Unix(),
	//}
	//var wg sync.WaitGroup
	//wg.Add(1)
	//err = a.sendEmail(message)
	//if err != nil {
	//	return nil, err
	//}
	//wg.Wait()

	return &domain.DefaultSigningResponse{
		AuthToken: response,
	}, nil
}

func (a authAppHandler) processLoginPasswordValidation(request domain.SigningRequest, identity *models.Identity) error {

	for _, credential := range identity.Credentials {
		if credential.Type == models.CredentialsTypeLogin {
			if credential.Status == models.CredentialStatusActive {

				err := a.encryptor.ComparePasscode(request.Password, credential.Password)
				if err != nil {
					a.logger.Error("SignIn", zap.Any(leetError.ErrorType(leetError.CredentialsValidationError), err), zap.Error(errors.New("credential password is not valid")))
					return leetError.ErrorResponseBody(leetError.CredentialsValidationError, err)
				}

				return nil
			}

			a.logger.Error("SignIn", zap.Error(leetError.ErrorResponseBody(leetError.UserLockedError, errors.New("credential status is not active"))))
			return leetError.ErrorResponseBody(leetError.UserLockedError, errors.New("credential status is not active"))
		}
	}

	a.logger.Error("SignIn", zap.Error(leetError.ErrorResponseBody(leetError.CredentialsValidationError, errors.New("credential type is not login"))))
	return leetError.ErrorResponseBody(leetError.CredentialsValidationError, errors.New("credential type is not login"))
}

func (a authAppHandler) resetPassword(customerID, firstName, lastName, email, passcode string, userCategory models.UserCategory) (*domain.DefaultSigningResponse, error) {

	hashedPasscode, err := a.passwordValidationEncryption(passcode)
	if err != nil {
		return nil, err
	}

	err = a.allRepository.AuthRepository.UpdateCredential(customerID, hashedPasscode)
	if err != nil {
		a.logger.Error("resetPassword", zap.Any("UpdateCredential", err))
		return nil, err
	}

	identity, err := a.allRepository.AuthRepository.GetIdentityByCustomerID(customerID)
	if err != nil {
		a.logger.Error("resetPassword", zap.Any("GetIdentityByCustomerID", err))
		return nil, err
	}

	response, err := a.tokenHandler.BuildAuthResponse(email, customerID, identity.ID, userCategory)
	if err != nil {
		a.logger.Error("SignIn", zap.Any("BuildAuthResponse", leetError.ErrorResponseBody(leetError.TokenGenerationError, err)))
		return nil, leetError.ErrorResponseBody(leetError.TokenGenerationError, err)
	}

	//var wg sync.WaitGroup
	//wg.Add(1)
	//message := models.Message{
	//	ID:         a.idGenerator.Generate(),
	//	Target:     email,
	//	TemplateID: library.ResetPasswordEmailTemplateID,
	//	DataMap: map[string]string{
	//		"FirstName": firstName,
	//		"LastName":  lastName,
	//	},
	//	Ts: time.Now().Unix(),
	//}
	//err = a.sendEmail(message)
	//if err != nil {
	//	return nil, err
	//}
	//wg.Wait()

	return &domain.DefaultSigningResponse{AuthToken: response}, nil
}

func (a authAppHandler) prepEmail(message models.Message, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	err := a.EmailClient.SendEmailWithTemplate(message)
	if err != nil {
		a.logger.Error("sendEmail", zap.Error(leetError.ErrorResponseBody(leetError.EmailSendingError, err)))
		errChan <- err
	}
}

func (a authAppHandler) sendEmail(message models.Message) error {
	var prepWg sync.WaitGroup
	//defer wg.Done()

	errChan := make(chan error, 1) // Use a buffered channel with a buffer size of 1
	prepWg.Add(1)
	go a.prepEmail(message, &prepWg, errChan)
	prepWg.Wait()

	select {
	case err := <-errChan:
		a.logger.Error("sendEmail", zap.Error(leetError.ErrorResponseBody(leetError.EmailSendingError, err)))
		return err
	default:
		return nil
	}
}
