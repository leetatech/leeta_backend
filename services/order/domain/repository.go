package domain

import (
	"context"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/leetatech/leeta_backend/services/models"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, request models.Order) error
	UpdateOrderStatus(ctx context.Context, request PersistOrderUpdate) error
	GetOrderByID(ctx context.Context, id string) (*models.Order, error)
	GetCustomerOrdersByStatus(ctx context.Context, request GetCustomerOrders) ([]OrderResponse, error)
	ListOrders(ctx context.Context, request query.ResultSelector, userId string) (orders []models.Order, totalResults uint64, err error)
	ListOrderStatusHistory(ctx context.Context, orderId string) ([]models.StatusHistory, error)
}
