package domain

import (
	"github.com/leetatech/leeta_backend/services/library/models"
)

type ProductRepository interface {
	CreateProduct(request models.Product) error
	GetProductByID(id string) (*models.Product, error)
}
