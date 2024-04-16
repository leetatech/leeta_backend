package domain

import (
	"context"
	"github.com/leetatech/leeta_backend/services/models"
)

type AddressRepository interface {
	Upsert(ctx context.Context, state models.State) error
	Update(ctx context.Context, state models.State) error
	GetState(ctx context.Context, name string) (models.State, error)
	GetAllStates(ctx context.Context) ([]models.State, error)
}
