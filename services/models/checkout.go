package models

import (
	"errors"
	"github.com/leetatech/leeta_backend/pkg/leetError"
)

type Checkout struct {
	ID              string          `json:"id" bson:"id"`
	CustomerID      string          `json:"customer_id" bson:"customer_id"`
	CheckoutDetails CheckoutDetails `json:"checkout_details" bson:"checkout_details"`
	ShippingInfo    ShippingInfo    `json:"shipping_info,omitempty" bson:"shipping_info"`
	AmountPaid      float64         `json:"amount_paid" bson:"amount_paid"`
	DeliveryFee     float64         `json:"delivery_fee" bson:"delivery_fee"`
	ServiceFee      float64         `json:"service_fee" bson:"service_fee"`
	TotalCost       float64         `json:"total_cost" bson:"total_cost"`
	Status          CheckoutStatus  `json:"status" bson:"status"`
	StatusTs        int64           `json:"status_ts" bson:"status_ts"`
	Ts              int64           `json:"ts" bson:"ts"`
} // @name Checkout

// ShippingInfo This object is only required from a guest/customer who is ordering for someone else
// If the guest/customer is ordering for himself, then we need to collect their address
type ShippingInfo struct {
	ForMe            bool    `json:"for_me" bson:"for_me"`
	RecipientName    string  `json:"recipient_name,omitempty" bson:"recipient_name"`
	RecipientPhone   string  `json:"recipient_phone,omitempty" bson:"recipient_phone"`
	RecipientEmail   string  `json:"recipient_email,omitempty" bson:"recipient_email"`
	RecipientAddress Address `json:"recipient_address,omitempty" bson:"recipient_address"`
} // @name ShippingInfo

type CheckoutDetails struct {
	CartItems []CartItem     `json:"cart_items" bson:"cart_items"`
	Status    CheckoutStatus `json:"status" bson:"status"`
	StatusTs  int64          `json:"status_ts" bson:"status_ts"`
	Ts        int64          `json:"ts" bson:"ts"`
} // @name CheckoutDetails

type CheckoutStatus string // @name CheckoutStatus

const (
	CheckoutCancelled CheckoutStatus = "cancelled" // @name CheckoutCancelled
	CheckoutAccepted  CheckoutStatus = "accepted"  // @name CheckoutAccepted
	CheckoutRejected  CheckoutStatus = "rejected"  // @name CheckoutRejected
	CheckoutPending   CheckoutStatus = "pending"   // @name CheckoutPending
	CheckoutFulFilled CheckoutStatus = "fulfilled" // @name CheckoutFulFilled
)

func IsValidCheckoutStatus(status CheckoutStatus) bool {
	return status == CheckoutCancelled || status == CheckoutAccepted || status == CheckoutRejected || status == CheckoutPending || status == CheckoutFulFilled
}

func SetCheckoutStatus(status CheckoutStatus) (CheckoutStatus, error) {
	switch IsValidCheckoutStatus(status) {
	case true:
		return status, nil
	default:
		return "", leetError.ErrorResponseBody(leetError.CheckoutStatusError, errors.New("invalid status provided"))
	}
}
