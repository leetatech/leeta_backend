package models

import (
	"fmt"
)

type Cart struct {
	ID         string       `json:"id" bson:"id"`
	CustomerID string       `json:"customer_id" bson:"customer_id"`
	CartItems  []CartItem   `json:"cart_items" bson:"cart_items"`
	Total      float64      `json:"total" bson:"total"`
	Status     CartStatuses `json:"status" bson:"status"`
	StatusTs   int64        `json:"status_ts" bson:"status_ts"`
	Ts         int64        `json:"ts" bson:"ts"`
}

type CartItem struct {
	ID              string          `json:"id" bson:"id"`
	ProductID       string          `json:"product_id" bson:"product_id"`
	ProductCategory ProductCategory `json:"product_category" bson:"product_category"`
	VendorID        string          `json:"vendor_id" bson:"vendor_id"`
	Weight          float32         `json:"weight,omitempty" bson:"weight"`
	Quantity        int             `json:"quantity,omitempty" bson:"quantity"`
	Cost            float64         `json:"cost" bson:"cost"`
}

func (c *CartItem) CalculateCartItemFee(fee *Fee) (float64, error) {
	var totalCost float64

	// Check if the product IDs match
	if fee.ProductID != c.ProductID {
		return 0, fmt.Errorf("cart product id: %s does not match fee's product id: %s", c.ProductID, fee.ProductID)
	}

	// check for quantity
	if c.Quantity == 0 {
		return 0, fmt.Errorf("invalid cart item, cart quantity cannot be zero %d", c.Quantity)
	}

	// Calculate cost based on weight or quantity
	if c.Weight != 0 {
		totalCost = float64(c.Weight) * fee.CostPerKg
	} else {
		totalCost = fee.CostPerQty
	}

	// Multiply cost by quantity
	totalCost *= float64(c.Quantity)

	return totalCost, nil
}

func (c *Cart) CalculateCartTotalFee() float64 {
	var totalCost float64

	for _, cartItem := range c.CartItems {
		totalCost += cartItem.Cost
	}

	return totalCost
}

type CartStatuses string

const (
	CartActive   CartStatuses = "ACTIVE"   // cart has been created and active
	CartInactive CartStatuses = "INACTIVE" // cart has been inactivated and no longer active due to check out or session expiry
)
