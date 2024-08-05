package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var invalidAppErr = errors.New("you are on the wrong app")

func (a authAppHandler) sendAccountVerificationEmail(ctx context.Context, fullName, userID, target, templateAlias string) error {
	requestOTP := domain.OTPRequest{
		Topic:  "Sign Up",
		Type:   models.EMAIL,
		Target: target,
	}
	otpResponse, err := a.createOTP(ctx, requestOTP)
	if err != nil {
		return fmt.Errorf("error creating OTP: %w", err)
	}
	err = a.notification.mail.Send(templateAlias, models.Message{
		ID:         a.idGenerator.Generate(),
		UserID:     userID,
		TemplateID: templateAlias,
		Title:      "Sign Up Verification",
		Sender:     a.mailerConfig.VerificationEmail,
		DataMap: map[string]string{
			"User": fullName,
			"OTP":  otpResponse.Message,
		},
		Recipients: []string{
			target,
		},
		Ts: time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("error sending verification email: %w", err)
	}

	return nil
}

func (a authAppHandler) validateAndEncryptPassword(password string) (string, error) {
	err := a.encryptor.ValidatePasswordStrength(password)
	if err != nil {
		return "", fmt.Errorf("error validating encryption password strength: %w", err)
	}
	passByte, err := a.encryptor.Generate(password)
	if err != nil {
		return "", errs.Body(errs.EncryptionError, err)
	}

	return string(passByte), nil
}

func (a authAppHandler) vendorSignUP(ctx context.Context, request domain.SignupRequest) (*domain.DefaultSigningResponse, error) {
	_, err := a.repositoryManager.AuthRepository.VendorByEmail(ctx, request.Email)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			timestamp := time.Now().Unix()

			vendor := models.Vendor{
				User: models.User{
					ID: a.idGenerator.Generate(),
					Email: models.Email{
						Address: request.Email,
					},
					Status: models.SignedUp,
				},
				TimeStamps: models.TimeStamps{
					Time: timestamp,
				},
			}
			err = vendor.User.ExtractName(request.FullName)
			if err != nil {
				return nil, errs.Body(errs.MissingUserNames, err)
			}
			err = a.repositoryManager.AuthRepository.CreateUser(ctx, vendor)
			if err != nil {
				return nil, err
			}

			identity := models.Identity{
				ID:       a.idGenerator.Generate(),
				UserID:   vendor.ID,
				Role:     models.VendorCategory,
				DeviceID: request.DeviceID,
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
			err = a.repositoryManager.AuthRepository.CreateIdentity(ctx, identity)
			if err != nil {
				return nil, errs.Body(errs.InternalError, err)
			}

			response, err := a.jwtManager.BuildAuthResponse(request.Email, vendor.ID, request.DeviceID, request.UserType)
			if err != nil {
				return nil, errs.Body(errs.TokenGenerationError, fmt.Errorf("error building authentication response on vendor sign up: %w", err))
			}

			err = a.sendAccountVerificationEmail(ctx, request.FullName, vendor.ID, vendor.Email.Address, pkg.VerifySignUPTemplatePath)
			if err != nil {
				return nil, err
			}
			return &domain.DefaultSigningResponse{AuthToken: response, Body: vendor.User}, nil

		default:
			return nil, errs.Body(errs.InternalError, err)
		}
	}

	return nil, errs.Body(errs.DuplicateUserError, errors.New("user already exists"))
}

func (a authAppHandler) customerSignUP(ctx context.Context, request domain.SignupRequest) (*domain.DefaultSigningResponse, error) {
	_, err := a.repositoryManager.AuthRepository.UserByEmail(ctx, request.Email)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			timestamp := time.Now().Unix()

			customer := models.Customer{
				User: models.User{
					ID: a.idGenerator.Generate(),
					Email: models.Email{
						Address: request.Email,
					},
					Status: models.SignedUp,
				},
				TimeStamps: models.TimeStamps{
					Time: timestamp,
				},
			}

			err = customer.User.ExtractName(request.FullName)
			if err != nil {
				return nil, errs.Body(errs.MissingUserNames, err)
			}

			err = a.repositoryManager.AuthRepository.CreateUser(ctx, customer)
			if err != nil {
				return nil, err
			}

			identity := models.Identity{
				ID:       a.idGenerator.Generate(),
				UserID:   customer.ID,
				Role:     models.CustomerCategory,
				DeviceID: request.DeviceID,
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
			err = a.repositoryManager.AuthRepository.CreateIdentity(ctx, identity)
			if err != nil {
				return nil, err
			}

			response, err := a.jwtManager.BuildAuthResponse(request.Email, customer.ID, request.DeviceID, request.UserType)
			if err != nil {
				return nil, errs.Body(errs.TokenGenerationError, fmt.Errorf("error building authentication response on customer sign up: %w", err))
			}

			err = a.sendAccountVerificationEmail(ctx, request.FullName, customer.ID, customer.Email.Address, pkg.VerifySignUPTemplatePath)
			if err != nil {
				return nil, errs.Body(errs.InternalError, err)
			}

			return &domain.DefaultSigningResponse{AuthToken: response, Body: customer.User}, nil

		default:
			return nil, errs.Body(errs.InternalError, err)
		}
	}

	return nil, errs.Body(errs.DuplicateUserError, errors.New("user already exists"))
}

func (a authAppHandler) buildSignIn(ctx context.Context, user models.User, status models.Statuses, request domain.SigningRequest) (*domain.DefaultSigningResponse, error) {
	identity, err := a.repositoryManager.AuthRepository.IdentityByUserID(ctx, user.ID)
	if err != nil {
		return nil, errs.Body(errs.IdentityNotFoundError, fmt.Errorf("error getting user identity by id %s when building sign in object: %w", user.ID, err))
	}

	err = a.validateLoginPassword(request, identity)
	if err != nil {
		return nil, err
	}

	switch status {
	case models.Locked, models.Exited, models.Rejected:
		return nil, errs.Body(errs.UserLockedError, fmt.Errorf("user with id %s is locked", user.ID))
	}

	response, err := a.jwtManager.BuildAuthResponse(request.Email, user.ID, request.DeviceID, request.UserType)
	if err != nil {
		return nil, errs.Body(errs.TokenGenerationError, fmt.Errorf("error building authentication response on sign in object: %w", err))
	}
	return &domain.DefaultSigningResponse{
		AuthToken: response,
		Body:      user,
	}, nil

}

func (a authAppHandler) vendorSignIN(ctx context.Context, request domain.SigningRequest) (*domain.DefaultSigningResponse, error) {
	vendor, err := a.repositoryManager.AuthRepository.VendorByEmail(ctx, request.Email)
	if err != nil {
		return nil, errs.Body(errs.UserNotFoundError, fmt.Errorf("error getting vendor identity by email %s when signing in: %w", request.Email, err))
	}

	if validateErr := a.validateUserRole(ctx, &request, &vendor.User); validateErr != nil {
		return nil, errs.Body(errs.InvalidUserRoleError, err)
	}

	return a.buildSignIn(ctx, vendor.User, vendor.Status, request)
}

func (a authAppHandler) customerSignIN(ctx context.Context, request domain.SigningRequest) (*domain.DefaultSigningResponse, error) {
	customer, err := a.repositoryManager.AuthRepository.UserByEmail(ctx, request.Email)
	if err != nil {
		return nil, errs.Body(errs.UserNotFoundError, fmt.Errorf("error getting customer identity by email %s when signing in: %w", request.Email, err))
	}

	if validateErr := a.validateUserRole(ctx, &request, &customer.User); validateErr != nil {
		return nil, errs.Body(errs.InvalidUserRoleError, validateErr)
	}

	return a.buildSignIn(ctx, customer.User, customer.Status, request)
}

func (a authAppHandler) validateUserRole(ctx context.Context, request *domain.SigningRequest, user *models.User) error {
	identity, err := a.repositoryManager.AuthRepository.IdentityByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	if identity.Role != request.UserType {
		return invalidAppErr
	}

	return nil
}

func (a authAppHandler) validateLoginPassword(request domain.SigningRequest, identity *models.Identity) error {
	for _, credential := range identity.Credentials {
		if credential.Type == models.CredentialsTypeLogin {
			if credential.Status == models.CredentialStatusActive {
				err := a.encryptor.ComparePasscode(request.Password, credential.Password)
				if err != nil {
					return errs.Body(errs.CredentialsValidationError, fmt.Errorf("error comparing passwords for validation %w", err))
				}
				return nil
			}
			return errs.Body(errs.UserLockedError, errors.New("credential status is not active"))
		}
	}
	return errs.Body(errs.CredentialsValidationError, errors.New("credential type is not login"))
}

func (a authAppHandler) createNewPassword(ctx context.Context, userID, passcode string) (*domain.DefaultSigningResponse, error) {
	hashedPasscode, err := a.validateAndEncryptPassword(passcode)
	if err != nil {
		return nil, err
	}

	err = a.repositoryManager.AuthRepository.UpdateCredential(ctx, userID, hashedPasscode)
	if err != nil {
		return nil, fmt.Errorf("error updating credentials: %w", err)
	}

	return &domain.DefaultSigningResponse{Body: "password reset successful"}, nil
}

func (a authAppHandler) adminSignUp(ctx context.Context, request domain.AdminSignUpRequest) (*domain.DefaultSigningResponse, error) {
	_, err := a.repositoryManager.AuthRepository.AdminByEmail(ctx, request.Email)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			timestamp := time.Now().Unix()

			admin := models.Admin{
				User: models.User{
					ID:        a.idGenerator.Generate(),
					FirstName: request.FirstName,
					LastName:  request.LastName,
					Email: models.Email{
						Address: request.Email,
					},
					DOB: request.DOB,
					Phone: models.Phone{
						Primary: true,
						Number:  request.Phone,
					},
					Status: models.SignedUp,
				},
				Department: request.Department,
				Role:       request.Role,
				TimeStamps: models.TimeStamps{
					Time: timestamp,
				},
			}

			admin.User.Address = append(admin.User.Address, models.Address{
				State:           request.Address.State,
				City:            request.Address.City,
				LGA:             request.Address.LGA,
				FullAddress:     request.Address.FullAddress,
				ClosestLandmark: request.Address.ClosestLandmark,
				AddressType:     models.CustomerResidentAddress,
			})

			err = a.repositoryManager.AuthRepository.CreateUser(ctx, admin)
			if err != nil {
				return nil, err
			}

			identity := models.Identity{
				ID:       a.idGenerator.Generate(),
				UserID:   admin.ID,
				Role:     models.AdminCategory,
				DeviceID: request.DeviceID,
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
			err = a.repositoryManager.AuthRepository.CreateIdentity(ctx, identity)
			if err != nil {
				return nil, err
			}

			response, err := a.jwtManager.BuildAuthResponse(request.Email, admin.ID, request.DeviceID, models.AdminCategory)
			if err != nil {
				return nil, errs.Body(errs.TokenGenerationError, fmt.Errorf("error building authentication response on admin sign up: %w", err))
			}

			err = a.sendAccountVerificationEmail(ctx, fmt.Sprintf("%s %s", request.FirstName, request.LastName), admin.ID, admin.User.Email.Address, pkg.AdminSignUpTemplatePath)
			if err != nil {
				return nil, err
			}

			return &domain.DefaultSigningResponse{AuthToken: response, Body: admin.User}, nil

		default:
			return nil, err
		}
	}

	return nil, errs.Body(errs.DuplicateUserError, nil)
}

func (a authAppHandler) adminSignIN(ctx context.Context, request domain.SigningRequest) (*domain.DefaultSigningResponse, error) {
	admin, err := a.repositoryManager.AuthRepository.AdminByEmail(ctx, request.Email)
	if err != nil {
		return nil, errs.Body(errs.UserNotFoundError, fmt.Errorf("error finding admin by email %s on admin sign in: %w", request.Email, err))
	}

	return a.buildSignIn(ctx, admin.User, admin.Status, request)
}
