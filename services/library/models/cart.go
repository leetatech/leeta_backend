package models

import (
	"errors"
	"github.com/leetatech/leeta_backend/services/library/leetError"
)

type Cart struct {
	ID          string       `json:"id" bson:"id"`
	CustomerID  string       `json:"customer_id" bson:"customer_id"`
	CartItems   []CartItem   `json:"cart_items" bson:"cart_items"`
	DeliveryFee float64      `json:"delivery_fee" bson:"delivery_fee"`
	Total       float64      `json:"total" bson:"total"`
	Status      CartStatuses `json:"status" bson:"status"`
	StatusTs    int64        `json:"status_ts" bson:"status_ts"`
	Ts          int64        `json:"ts" bson:"ts"`
}

type CartItem struct {
	ID         string  `json:"id" bson:"id"`
	CustomerID string  `json:"customer_id" bson:"customer_id"`
	SessionID  string  `json:"session_id" bson:"session_id"`
	ProductID  string  `json:"product_id" bson:"product_id"`
	VendorID   string  `json:"vendor_id" bson:"vendor_id"`
	Weight     float32 `json:"weight" bson:"weight"`
	Quantity   int32   `json:"quantity" bson:"quantity"`
	AmountPaid float64 `json:"amount_paid" bson:"amount_paid"`
}

type CartStatuses string

const (
	CartActive   CartStatuses = "ACTIVE"   // order has been created and processing
	CartInactive CartStatuses = "INACTIVE" // order was rejected by vendor or customer
)

func IsValidCartStatus(status CartStatuses) bool {
	return status == CartActive || status == CartInactive
}

func SetCartStatus(status CartStatuses) (CartStatuses, error) {
	switch IsValidCartStatus(status) {
	case true:
		return status, nil
	default:
		return "", leetError.ErrorResponseBody(leetError.CartStatusesError, errors.New("invalid order status"))
	}
}
