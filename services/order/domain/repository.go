package domain

import (
	"context"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/leetatech/leeta_backend/services/models"
)

type OrderRepository interface {
	Create(ctx context.Context, request models.Order) error
	UpdateStatus(ctx context.Context, request PersistOrderUpdate) error
	OrderByID(ctx context.Context, id string) (*models.Order, error)
	OrdersByStatus(ctx context.Context, request GetCustomerOrders) ([]Response, error)
	Orders(ctx context.Context, request query.ResultSelector, userId string) (orders []models.Order, totalResults uint64, err error)
	OrderStatusHistory(ctx context.Context, orderId string) ([]models.StatusHistory, error)
}
