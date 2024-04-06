package domain

import "errors"

type InactivateCart struct {
	ID string `json:"id"`
} // @name InactivateCart

type CartItem struct {
	ProductID string  `json:"product_id" bson:"product_id"`
	Weight    float32 `json:"weight,omitempty" bson:"weight"`
	Quantity  int     `json:"quantity,omitempty" bson:"quantity"`
	Cost      float64 `json:"cost" bson:"cost"`
} // @name CartRefillDetails

type UpdateCartItemQuantityRequest struct {
	CartItemID string `json:"cart_item_id"`
	Quantity   int    `json:"quantity"`
} // @name UpdateCartItemQuantityRequest

func (u *UpdateCartItemQuantityRequest) IsValid() (bool, error) {
	if u.CartItemID == "" {
		return false, errors.New("cart_item_id is empty")
	}
	if u.Quantity <= 0 {
		return false, errors.New("quantity is invalid")
	}
	return true, nil
}
