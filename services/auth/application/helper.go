package application

import (
	"github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

func (a authAppHandler) passwordValidationEncryption(password string) (string, error) {
	err := a.encryptor.IsValidPassword(password)
	if err != nil {
		return "", err
	}
	passByte, err := a.encryptor.GenerateFromPasscode(password)
	if err != nil {
		return "", leetError.ErrorResponseBody(leetError.DecryptionError, err)
	}

	return string(passByte), nil
}

func (a authAppHandler) vendorSignUP(request domain.SignUpRequest) (*domain.DefaultSigningResponse, error) {
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

			response, err := a.tokenHandler.BuildAuthResponse(request.Email, vendor.ID, a.idGenerator.Generate())
			if err != nil {
				return nil, err
			}

			requestOTP := domain.OTPRequest{
				Topic:        "Sign Up",
				Type:         models.EMAIL,
				Target:       request.Email,
				UserCategory: models.VendorCategory,
			}
			_, err = a.CreateOTP(requestOTP)
			if err != nil {
				a.logger.Error("SignUp", zap.Any("CreateOTP", err))
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
