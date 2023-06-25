package application

import (
	"github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/library/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type authAppHandler struct {
	tokenHandler  library.TokenHandler
	encryptor     library.EncryptorManager
	idGenerator   library.IDGenerator
	logger        *zap.Logger
	allRepository library.Repositories
}

type AuthApplication interface {
	SignUp(request domain.SignUpRequest) (*domain.DefaultSigningResponse, error)
}

func NewAuthApplication(tokenHandler library.TokenHandler, logger *zap.Logger, allRepository library.Repositories) AuthApplication {
	return &authAppHandler{
		tokenHandler:  tokenHandler,
		encryptor:     library.NewEncryptor(),
		idGenerator:   library.NewIDGenerator(),
		allRepository: allRepository,
	}
}

func (a authAppHandler) SignUp(request domain.SignUpRequest) (*domain.DefaultSigningResponse, error) {
	err := a.encryptor.IsValidPassword(request.Password)
	if err != nil {
		return nil, err
	}
	if models.IsValidUserCategory(request.UserType) {
		switch request.UserType {
		case models.VendorCategory:
			return a.vendorSignUP(request)
		}
	}

	return nil, nil
}

func (a authAppHandler) vendorSignUP(request domain.SignUpRequest) (*domain.DefaultSigningResponse, error) {
	_, err := a.allRepository.AuthRepository.GetVendorByEmail(request.Email)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			break
		default:
			return nil, err
		}
	}

	hashedPasscode, err := a.encryptor.GenerateFromPasscode(request.Password)
	if err != nil {
		return nil, err
	}

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
				Password:        string(hashedPasscode),
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

	return &domain.DefaultSigningResponse{AuthToken: response}, nil
}
