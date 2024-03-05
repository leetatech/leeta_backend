package domain

import (
	"context"
	"github.com/leetatech/leeta_backend/services/library/models"
)

type FeesRepository interface {
	CreateFees(ctx context.Context, request models.Fees) error
	GetFeeByProductID(ctx context.Context, productID string, status models.FeesStatuses) (*models.Fees, error)
	GetFees(ctx context.Context, status models.FeesStatuses) ([]models.Fees, error)
	UpdateFees(ctx context.Context, status models.FeesStatuses) error
}
