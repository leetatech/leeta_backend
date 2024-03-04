package domain

type InactivateCart struct {
	ID string `json:"id"`
} // @name InactivateCart

type AddToCartRequest struct {
	Guest         bool     `json:"guest" bson:"guest"`
	RefillDetails CartItem `json:"refill_details" bson:"refill_details"`
} // @name AddToCartRequest

type CartItem struct {
	ProductID string  `json:"product_id" bson:"product_id"`
	Weight    float32 `json:"weight" bson:"weight"`
	Cost      float64 `json:"cost" bson:"cost"`
} // @name CartRefillDetails
