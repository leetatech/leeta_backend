package domain

import "github.com/leetatech/leeta_backend/services/models"

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

type UpdateCartItemQuantity struct {
	CartItemID    string  `json:"cart_item_id"`
	Quantity      int     `json:"quantity"`
	ItemTotalCost float64 `json:"item_total_cost"`
	CartTotalCost float64 `json:"cart_total_cost"`
}

type StoredCartItemDetails struct {
	ProductID       string                 `json:"product_id" bson:"product_id"`
	ProductCategory models.ProductCategory `json:"product_category" bson:"product_category"`
	Weight          float32                `json:"weight,omitempty" bson:"weight"`
	Quantity        int                    `json:"quantity,omitempty" bson:"quantity"`
}
