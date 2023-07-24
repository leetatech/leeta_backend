package domain

import "github.com/leetatech/leeta_backend/services/library/models"

type UserRepository interface {
	VendorDetailsUpdate(request VendorDetailsUpdateRequest) error
	RegisterVendorBusiness(request models.Business) error
	GetVendorByID(id string) (*models.Vendor, error)
}
