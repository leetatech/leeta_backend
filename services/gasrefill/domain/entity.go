package domain

import "github.com/leetatech/leeta_backend/services/library/models"

type GasRefillRequest struct {
	Guest         bool                `json:"guest" bson:"guest"`
	GuestBioData  models.GuestBioData `json:"guest_bio_data,omitempty" bson:"guest_bio_data"`
	CustomerID    string              `json:"customer_id" bson:"customer_id"`
	RefillDetails RefillDetails       `json:"refill_details" bson:"refill_details"`
	// This object is to be sent when the customer is done with their order and payment
	ShippingInfo models.ShippingInfo `json:"shipping_info,omitempty" bson:"shipping_info"`
} // @name GasRefillRequest

type RefillDetails struct {
	ProductID  string                 `json:"product_id" bson:"product_id"`
	Weight     float32                `json:"weight" bson:"weight"`
	AmountPaid float64                `json:"amount_paid" bson:"amount_paid"`
	GasType    models.ProductCategory `json:"gas_type" bson:"gas_type"`
} // @name RefillDetails

type UpdateRefillRequest struct {
	RefillID      string                     `json:"refill_id" bson:"refill_id"`
	RefillDetails models.RefillDetails       `json:"refill_details" bson:"refill_details"`
	RequestStatus models.RefillRequestStatus `json:"request_status" bson:"request_status"`
	Reason        string                     `json:"reason" bson:"reason"`
} // @name UpdateRefillRequest

type ListRefillFilter struct {
	Status     []models.RefillRequestStatus `json:"status" bson:"status"`
	GasType    []models.ProductCategory     `json:"gas_type" bson:"gas_type"`
	CustomerID string                       `json:"customer_id" json:"customer_id"`
	GuestEmail string                       `json:"guest_email" bson:"guest_email"`
	Limit      int64                        `json:"limit" bson:"limit"`
	Page       int64                        `json:"page" bson:"page"`
} // @name ListRefillFilter
