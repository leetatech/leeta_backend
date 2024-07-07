package domain

import (
	"context"
	"github.com/leetatech/leeta_backend/services/models"
)

type AuthRepository interface {
	CreateIdentity(ctx context.Context, identity models.Identity) error
	CreateGuestRecord(ctx context.Context, guest models.Guest) error
	GetGuestRecord(ctx context.Context, deviceId string) (models.Guest, error)
	GetVendorByEmail(ctx context.Context, email string) (*models.Vendor, error)
	CreateOTP(ctx context.Context, verifications models.Verification) error
	EarlyAccess(ctx context.Context, earlyAccess models.EarlyAccess) error
	GetIdentityByUserID(ctx context.Context, id string) (*models.Identity, error)
	GetOTPForValidation(ctx context.Context, target string) (*models.Verification, error)
	ValidateOTP(ctx context.Context, verificationId string) error
	UpdateCredential(ctx context.Context, customerID, password string) error
	GetAdminByEmail(ctx context.Context, email string) (*models.Admin, error)
	GetUserByEmail(ctx context.Context, email string) (*models.Customer, error)
	CreateUser(ctx context.Context, user any) error
	UpdateEmailVerify(ctx context.Context, email string, status bool) error
	UpdateGuestRecord(ctx context.Context, guest models.Guest) error
	GetUserByEmailOrPhone(ctx context.Context, target string) (*models.Customer, error)
	UpdatePhoneVerify(ctx context.Context, phone string, status bool) error
}
