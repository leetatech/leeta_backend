package domain

import (
	"context"
	"github.com/leetatech/leeta_backend/services/library/models"
)

type FeesRepository interface {
	CreateFees(ctx context.Context, request models.Fee) error
	GetFeeByProductID(ctx context.Context, productID string, status models.FeesStatuses) (*models.Fee, error)
	GetFees(ctx context.Context, status models.FeesStatuses) ([]models.Fee, error)
	UpdateFees(ctx context.Context, status models.FeesStatuses) error
}
