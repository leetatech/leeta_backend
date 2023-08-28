package models

import (
	"errors"
	"github.com/leetatech/leeta_backend/services/library/leetError"
)

type Order struct {
	ID          string        `json:"id" bson:"id"`
	ProductID   string        `json:"product_id" bson:"product_id"`
	CustomerID  string        `json:"customer_id" bson:"customer_id"`
	VendorID    string        `json:"vendor_id" bson:"vendor_id"`
	VAT         float64       `json:"vat" bson:"vat"`
	DeliveryFee float64       `json:"delivery_fee" bson:"delivery_fee"`
	Total       float64       `json:"total" bson:"total"`
	Status      OrderStatuses `json:"status" bson:"status"`
	StatusTs    int64         `json:"status_ts" bson:"status_ts"`
	Ts          int64         `json:"ts" bson:"ts"`
} // @name Order

// OrderStatuses type
type OrderStatuses string

const (
	OrderPending   OrderStatuses = "PENDING"   // order has been created and processing
	OrderRejected  OrderStatuses = "REJECTED"  // order has been rejected by vendor or customer
	OrderCompleted OrderStatuses = "COMPLETED" // order has been processed and delivered
	OrderApproved  OrderStatuses = "APPROVED"  // order has been processed, delivered, and verified by the customer
)

func IsValidOrderStatus(status OrderStatuses) bool {
	return status == OrderPending || status == OrderRejected || status == OrderCompleted || status == OrderApproved
}

func SetOrderStatus(status OrderStatuses) (OrderStatuses, error) {
	switch IsValidOrderStatus(status) {
	case true:
		return status, nil
	default:
		return "", leetError.ErrorResponseBody(leetError.OrderStatusesError, errors.New("invalid order status"))
	}
}
