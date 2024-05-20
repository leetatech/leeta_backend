package models

import (
	"errors"
	"github.com/leetatech/leeta_backend/pkg/leetError"
)

type Order struct {
	ID              string       `json:"id" bson:"id"`
	Orders          []CartItem   `json:"orders" bson:"orders"`
	CustomerID      string       `json:"customer_id" bson:"customer_id"`
	DeliveryDetails ShippingInfo `json:"delivery_details" bson:"delivery_details"`
	PaymentMethod   string       `json:"payment_method" bson:"payment_method"`
	//VendorID        string        `json:"vendor_id" bson:"vendor_id"` // uncomment vendor id when sure how vendors affects individual orders
	DeliveryFee float64       `json:"delivery_fee" bson:"delivery_fee"`
	ServiceFee  float64       `json:"service_fee" bson:"service_fee"`
	Total       float64       `json:"total" bson:"total"`
	Status      OrderStatuses `json:"status" bson:"status"`
	StatusTs    int64         `json:"status_ts" bson:"status_ts"`
	Ts          int64         `json:"ts" bson:"ts"`
} // @name Order

// ShippingInfo is the object required for shipping details of an order
type ShippingInfo struct {
	Name    string  `json:"name,omitempty" bson:"name"`
	Phone   string  `json:"phone,omitempty" bson:"phone"`
	Email   string  `json:"email,omitempty" bson:"email"`
	Address Address `json:"address,omitempty" bson:"address"`
} // @name ShippingInfo

// OrderStatuses type
type OrderStatuses string

const (
	OrderPending   OrderStatuses = "PENDING"   // order has been created and processing
	OrderCancelled OrderStatuses = "CANCELLED" // order has been cancelled by vendor or customer
	OrderRejected  OrderStatuses = "REJECTED"  // order was rejected by vendor or customer
	OrderCompleted OrderStatuses = "COMPLETED" // order has been processed and delivered, and verified by the customer
	OrderApproved  OrderStatuses = "APPROVED"  // order has been approved
)

func IsValidOrderStatus(status OrderStatuses) bool {
	return status == OrderPending || status == OrderCancelled || status == OrderRejected || status == OrderCompleted || status == OrderApproved
}

func SetOrderStatus(status OrderStatuses) (OrderStatuses, error) {
	switch IsValidOrderStatus(status) {
	case true:
		return status, nil
	default:
		return "", leetError.ErrorResponseBody(leetError.OrderStatusesError, errors.New("invalid order status"))
	}
}
