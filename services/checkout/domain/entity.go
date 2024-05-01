package domain

import "github.com/leetatech/leeta_backend/services/models"

type CheckoutRequest struct {
	CartID          string              `json:"cart_id" bson:"cart_id"`
	DeliveryDetails models.ShippingInfo `json:"delivery_details" bson:"delivery_details"`
	PaymentMethod   string              `json:"payment_method" bson:"payment_method"`
	AmountPaid      float64             `json:"amount_paid" bson:"amount_paid"`
	DeliveryFee     float64             `json:"delivery_fee" bson:"delivery_fee"`
	ServiceFee      float64             `json:"service_fee" bson:"service_fee"`
} // @name CheckoutRequest

type UpdateRefillRequest struct {
	RefillID      string                `json:"refill_id" bson:"refill_id"`
	RequestStatus models.CheckoutStatus `json:"request_status" bson:"request_status"`
	Reason        string                `json:"reason" bson:"reason"`
} // @name UpdateRefillRequest

type ListRefillFilter struct {
	Status     []models.CheckoutStatus  `json:"status" bson:"status"`
	GasType    []models.ProductCategory `json:"gas_type" bson:"gas_type"`
	CustomerID string                   `json:"customer_id" bson:"customer_id"`
	GuestEmail string                   `json:"guest_email" bson:"guest_email"`
	Limit      int64                    `json:"limit" bson:"limit"`
	Page       int64                    `json:"page" bson:"page"`
} // @name ListRefillFilter
