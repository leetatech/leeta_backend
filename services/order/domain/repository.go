package domain

import "github.com/leetatech/leeta_backend/services/library/models"

type OrderRepository interface {
	CreateOrder(request OrderRequest) (*OrderResponse, error)
	UpdateOrderStatus(request models.OrderStatuses) error
}
