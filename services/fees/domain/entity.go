package domain

type FeeQuotationRequest struct {
	CostPerKg  float64 `json:"cost_per_kg,omitempty" bson:"cost_per_kg"`
	CostPerQty float64 `json:"cost_per_qty,omitempty" bson:"cost_per_qty"`
	ProductID  string  `json:"product_id" bson:"product_id"`
} // @name FeeQuotationRequest
