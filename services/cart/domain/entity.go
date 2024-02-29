package domain

import "github.com/leetatech/leeta_backend/services/library/models"

type AddToCart struct {
	CartItems   []models.CartItem `json:"cart_items"`
	DeliveryFee float64           `json:"delivery_fee" bson:"delivery_fee"`
	Total       float64           `json:"total" bson:"total"`
} // @name AddToCart

type InactivateCart struct {
	ID string `json:"id"`
} // @name InactivateCart
