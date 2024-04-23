package domain

import (
	"context"
	"github.com/leetatech/leeta_backend/pkg/query"
	"github.com/leetatech/leeta_backend/services/models"
)

type FeesRepository interface {
	CreateFees(ctx context.Context, request models.Fee) error
	GetFeeByProductID(ctx context.Context, productID string, status models.FeesStatuses) (*models.Fee, error)
	GetFeesByStatus(ctx context.Context, status models.FeesStatuses) ([]models.Fee, error)
	UpdateFees(ctx context.Context, status models.FeesStatuses, feeType models.FeeType, lga models.LGA, productID string) error
	GetTypedFees(ctx context.Context, request query.ResultSelector) ([]models.Fee, uint64, error)
}
