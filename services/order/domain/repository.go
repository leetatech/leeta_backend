package domain

import "github.com/leetatech/leeta_backend/services/library/models"

type OrderRepository interface {
	CreateOrder(request OrderRequest) error
	UpdateOrderStatus(request UpdateOrderStatusRequest) error
	GetOrderByID(id string) (*models.Order, error)
}
