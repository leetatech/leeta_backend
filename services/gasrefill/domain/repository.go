package domain

import (
	"context"
	"github.com/leetatech/leeta_backend/services/library/models"
)

type GasRefillRepository interface {
	RequestRefill(ctx context.Context, request models.GasRefill) error
	CreateFees(ctx context.Context, request models.Fees) error
	GetFees(ctx context.Context, status models.CartStatuses) (*models.Fees, error)
	UpdateFees(ctx context.Context, status models.CartStatuses) error
}
