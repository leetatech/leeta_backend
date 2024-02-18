package application

import (
	"context"
	"github.com/leetatech/leeta_backend/services/gasrefill/domain"
	"github.com/leetatech/leeta_backend/services/library"
)

type GasRefillHandler struct {
}

type GasRefillApplication interface {
	RefillGas(ctx context.Context, request domain.GasRefillRequest) (*library.DefaultResponse, error)
	AcceptRefill(ctx context.Context, request domain.UpdateRefillRequest) (*library.DefaultResponse, error)
	CancelRefill(ctx context.Context, request domain.UpdateRefillRequest) (*library.DefaultResponse, error)
	RejectRefill(ctx context.Context, request domain.UpdateRefillRequest) (*library.DefaultResponse, error)
	CompleteRefill(ctx context.Context, refillID string) (*library.DefaultResponse, error)
}

func NewGasRefillApplication(request library.DefaultApplicationRequest) GasRefillApplication {
	return &GasRefillHandler{}
}

func (r *GasRefillHandler) RefillGas(ctx context.Context, request domain.GasRefillRequest) (*library.DefaultResponse, error) {
	return nil, nil
}

func (r *GasRefillHandler) AcceptRefill(ctx context.Context, request domain.UpdateRefillRequest) (*library.DefaultResponse, error) {
	return nil, nil
}

func (r *GasRefillHandler) CancelRefill(ctx context.Context, request domain.UpdateRefillRequest) (*library.DefaultResponse, error) {
	return nil, nil
}

func (r *GasRefillHandler) RejectRefill(ctx context.Context, request domain.UpdateRefillRequest) (*library.DefaultResponse, error) {
	return nil, nil
}

func (r *GasRefillHandler) CompleteRefill(ctx context.Context, refillID string) (*library.DefaultResponse, error) {
	return nil, nil
}
