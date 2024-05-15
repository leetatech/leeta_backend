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

type CartCheckoutRequest struct {
	CartID          string              `json:"cart_id" bson:"cart_id"`
	DeliveryDetails models.ShippingInfo `json:"delivery_details" bson:"delivery_details"`
	PaymentMethod   string              `json:"payment_method" bson:"payment_method"`
	DeliveryFee     float64             `json:"delivery_fee" bson:"delivery_fee"`
	ServiceFee      float64             `json:"service_fee" bson:"service_fee"`
	TotalFee        float64             `json:"total_fee" bson:"total_fee"`
} // @name CartCheckoutRequest
