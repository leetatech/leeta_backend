package domain

type InactivateCart struct {
	ID string `json:"id"`
} // @name InactivateCart

type AddToCartRequest struct {
	Guest       bool     `json:"guest" bson:"guest"`
	CartDetails CartItem `json:"cart_details" bson:"cart_details"`
} // @name AddToCartRequest

type CartItem struct {
	ProductID string  `json:"product_id" bson:"product_id"`
	Weight    float32 `json:"weight,omitempty" bson:"weight"`
	Quantity  int     `json:"quantity,omitempty" bson:"quantity"`
	Cost      float64 `json:"cost" bson:"cost"`
} // @name CartRefillDetails

type DeleteCartItemRequest struct {
	CartItemID string `json:"cart_item_id"`
} // @name DeleteCartItemRequest
