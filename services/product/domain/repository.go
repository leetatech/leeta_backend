package domain

import (
	"context"
	"github.com/leetatech/leeta_backend/pkg/query"
	"github.com/leetatech/leeta_backend/services/models"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, request models.Product) error
	GetProductByID(ctx context.Context, id string) (models.Product, error)
	GetAllVendorProducts(ctx context.Context, request GetVendorProductsRequest) ([]models.Product, error)
	ListProducts(ctx context.Context, request query.ResultSelector) (products []models.Product, totalResults uint64, err error)
}
