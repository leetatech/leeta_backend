package domain

import (
	"context"
	"github.com/leetatech/leeta_backend/services/models"
)

type CheckoutRepository interface {
	RequestCheckout(ctx context.Context, request models.Checkout) error
	UpdateCheckoutStatus(ctx context.Context, id string, status models.CheckoutStatus) error
}
