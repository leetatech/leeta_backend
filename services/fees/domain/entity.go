package domain

type FeeQuotationRequest struct {
	CostPerKg  float64 `json:"cost_per_kg,omitempty" bson:"cost_per_kg"`
	CostPerQty float64 `json:"cost_per_qty,omitempty" bson:"cost_per_qty"`
	ServiceFee float64 `json:"service_fee" bson:"service_fee"`
	ProductID  string  `json:"product_id" bson:"product_id"`
} // @name FeeQuotationRequest
