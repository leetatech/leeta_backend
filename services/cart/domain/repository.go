package domain

import (
	"context"

	"github.com/leetatech/leeta_backend/services/models"
)

type CartRepository interface {
	AddToCart(ctx context.Context, request models.Cart) error
	GetCartByCustomerID(ctx context.Context, customerID string) (*models.Cart, error)
	GetCartByDeviceID(ctx context.Context, deviceID string) (*models.Cart, error)
	UpdateCart(ctx context.Context, request models.Cart) error
	AddToCartItem(ctx context.Context, cartID string, cartItems models.CartItem, total float64, statusTs int64) error
	DeleteCartItem(ctx context.Context, cartItemID string, itemTotalCost float64) error
	InactivateCart(ctx context.Context, id string) error
	GetCartByCartItemID(ctx context.Context, cartItemID string) (*models.Cart, error)
}
