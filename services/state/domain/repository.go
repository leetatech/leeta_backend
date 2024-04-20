package domain

import (
	"context"
	"github.com/leetatech/leeta_backend/services/models"
)

type StateRepository interface {
	SaveStates(ctx context.Context, states []any) error
	GetState(ctx context.Context, name string) (models.State, error)
	GetAllStates(ctx context.Context) ([]models.State, error)
}
