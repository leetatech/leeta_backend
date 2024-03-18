package domain

import (
	"context"
	"github.com/leetatech/leeta_backend/services/models"
)

type GasRefillRepository interface {
	RequestRefill(ctx context.Context, request models.GasRefill) error
}
