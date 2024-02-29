package domain

import "github.com/leetatech/leeta_backend/services/library/models"

type InactivateCart struct {
	ID       string `json:"id"`
	DeviceID string `json:"device_id"`
} // @name InactivateCart

type AddToCartRequest struct {
	Guest         bool              `json:"guest" bson:"guest"`
	RefillDetails CartRefillDetails `json:"refill_details" bson:"refill_details"`
} // @name AddToCartRequest

type CartRefillDetails struct {
	ProductID string                 `json:"product_id" bson:"product_id"`
	Weight    float32                `json:"weight" bson:"weight"`
	CostPerKg float64                `json:"cost_per_kg" bson:"cost_per_kg"`
	GasType   models.ProductCategory `json:"gas_type" bson:"gas_type"`
} // @name CartRefillDetails
