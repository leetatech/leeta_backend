package domain

import (
	"errors"
	"github.com/leetatech/leeta_backend/services/models"
)

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

type ListCartResponse struct {
	Cart        CartResponse `json:"cart"`
	HasNextPage bool         `json:"has_next_page"`
} // @name ListCartResponse

type CartResponse struct {
	ID           string            `json:"id" bson:"id"`
	CartItems    []models.CartItem `json:"cart_items" bson:"cart_items"`
	Total        float64           `json:"total" bson:"total"`
	TotalRecords int               `json:"total_records" bson:"total_records"`
} // @name CartResponse
