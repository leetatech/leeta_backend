package domain

import "github.com/leetatech/leeta_backend/services/library/models"

type Gas struct {
	ID              string              `json:"id" bson:"id"`
	CustomerID      string              `json:"customer_id" bson:"customer_id"`
	GasType         GasType             `json:"gas_type" bson:"gas_type"`
	Status          RefillRequestStatus `json:"status" bson:"status"`
	Weight          float32             `json:"weight" bson:"weight"`
	DeliveryAddress models.Address      `json:"delivery_address" bson:"delivery_address"`
	AmountPaid      float64             `json:"amount_paid "bson:"amount_paid"`
} // @name Gas

type GasRefillRequest struct {
	CustomerID      string         `json:"customer_id" bson:"customer_id"`
	GasType         GasType        `json:"gas_type" bson:"gas_type"`
	Weight          float32        `json:"weight" bson:"weight"`
	DeliveryAddress models.Address `json:"delivery_address" bson:"delivery_address"`
	AmountPaid      float64        `json:"amount_paid "bson:"amount_paid"`
} // @name GasRefillRequest

type GasType string

const (
	LPG GasType = "lpg"
	LNG GasType = "lng"
)

type UpdateRefillRequest struct {
	RefillID      string              `json:"refill_id" bson:"refill_id"`
	RequestStatus RefillRequestStatus `json:"request_status" bson:"request_status"`
	Reason        string              `json:"reason" bson:"reason"`
}

type RefillRequestStatus string

const (
	Cancelled RefillRequestStatus = "cancelled"
	Accepted  RefillRequestStatus = "accepted"
	Rejected  RefillRequestStatus = "rejected"
	Pending   RefillRequestStatus = "pending"
	FulFilled RefillRequestStatus = "fulfilled"
)

type ListRefillFilter struct {
	Status     []RefillRequestStatus `json:"status" bson:"status"`
	CustomerID string                `json:"customer_id" json:"customer_id"`
	Limit      int64                 `json:"limit" bson:"limit"`
	Page       int64                 `json:"page" bson:"page"`
} // @name ListRefillFilter
