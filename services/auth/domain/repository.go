package domain

import (
	"github.com/leetatech/leeta_backend/services/library/models"
)

type AuthRepository interface {
	CreateVendor(vendor models.Vendor) error
	CreateIdentity(identity models.Identity) error
	GetVendorByEmail(email string) (*models.Vendor, error)
}
