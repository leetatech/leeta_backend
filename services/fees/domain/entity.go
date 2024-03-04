package domain

type FeeQuotationRequest struct {
	CostPerKg  float64 `json:"cost_per_kg" bson:"cost_per_kg"`
	ServiceFee float64 `json:"service_fee" bson:"service_fee"`
	ProductID  string  `json:"product_id" bson:"product_id"`
} // @name FeeQuotationRequest
