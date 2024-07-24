package domain

import (
	"context"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/leetatech/leeta_backend/services/models"
)

type ProductRepository interface {
	Create(ctx context.Context, request models.Product) error
	Product(ctx context.Context, id string) (models.Product, error)
	VendorProducts(ctx context.Context, request GetVendorProductsRequest) ([]models.Product, error)
	ListProducts(ctx context.Context, request query.ResultSelector) (products []models.Product, totalResults uint64, err error)
}
