package domain

import "github.com/leetatech/leeta_backend/services/models"

type UserRepository interface {
	VendorDetailsUpdate(request VendorDetailsUpdateRequest) error
	RegisterVendorBusiness(request models.Business) error
	GetVendorByID(id string) (*models.Vendor, error)
	GetCustomerByID(id string) (*models.Customer, error)
}
