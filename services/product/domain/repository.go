package domain

import (
	"context"
	"github.com/leetatech/leeta_backend/services/models"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, request models.Product) error
	CreateGasProduct(ctx context.Context, request models.Product) error
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	GetAllVendorProducts(ctx context.Context, request GetVendorProductsRequest) (*GetVendorProductsResponse, error)
}
