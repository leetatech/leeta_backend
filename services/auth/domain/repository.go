package domain

import (
	"context"
	"github.com/leetatech/leeta_backend/services/library/models"
)

type AuthRepository interface {
	CreateVendor(ctx context.Context, vendor models.Vendor) error
	CreateIdentity(ctx context.Context, identity models.Identity) error
	GetVendorByEmail(ctx context.Context, email string) (*models.Vendor, error)
	CreateOTP(ctx context.Context, verifications models.Verification) error
	EarlyAccess(ctx context.Context, earlyAccess models.EarlyAccess) error
	GetIdentityByCustomerID(ctx context.Context, id string) (*models.Identity, error)
	GetOTPForValidation(ctx context.Context, target string) (*models.Verification, error)
	ValidateOTP(ctx context.Context, verificationId string) error
	UpdateCredential(ctx context.Context, customerID, password string) error
	GetAdminByEmail(ctx context.Context, email string) (*models.Admin, error)
	CreateAdmin(ctx context.Context, admin models.Admin) error
	CreateCustomer(ctx context.Context, customer models.Customer) error
	GetCustomerByEmail(ctx context.Context, email string) (*models.Customer, error)
}
