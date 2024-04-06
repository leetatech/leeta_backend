package domain

import (
	"context"
	"github.com/leetatech/leeta_backend/services/models"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, request models.Order) error
	UpdateOrderStatus(ctx context.Context, request UpdateOrderStatusRequest) error
	GetOrderByID(ctx context.Context, id string) (*models.Order, error)
	GetCustomerOrdersByStatus(ctx context.Context, request GetCustomerOrders) ([]OrderResponse, error)
}
