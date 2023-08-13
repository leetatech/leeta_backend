package domain

import (
	"github.com/leetatech/leeta_backend/services/library/models"
)

type AuthRepository interface {
	CreateVendor(vendor models.Vendor) error
	CreateIdentity(identity models.Identity) error
	GetVendorByEmail(email string) (*models.Vendor, error)
	CreateOTP(verifications models.Verification) error
	EarlyAccess(earlyAccess models.EarlyAccess) error
	GetIdentityByCustomerID(id string) (*models.Identity, error)
	GetOTPForValidation(target string) (*models.Verification, error)
	ValidateOTP(verificationId string) error
	UpdateCredential(customerID, password string) error
	GetAdminByEmail(email string) (*models.Admin, error)
	CreateAdmin(admin models.Admin) error
}
