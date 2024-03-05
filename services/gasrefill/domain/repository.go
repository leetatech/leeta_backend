package domain

import (
	"context"
	"github.com/leetatech/leeta_backend/services/library/models"
)

type GasRefillRepository interface {
	RequestRefill(ctx context.Context, request models.GasRefill) error
}
