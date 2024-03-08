package models

import (
	"errors"
	"github.com/leetatech/leeta_backend/services/library/leetError"
)

type Cart struct {
	ID         string       `json:"id" bson:"id"`
	CustomerID string       `json:"customer_id" bson:"customer_id"`
	DeviceID   string       `json:"device_id" bson:"device_id"`
	CartItems  []CartItem   `json:"cart_items" bson:"cart_items"`
	Total      float64      `json:"total" bson:"total"`
	Status     CartStatuses `json:"status" bson:"status"`
	StatusTs   int64        `json:"status_ts" bson:"status_ts"`
	Ts         int64        `json:"ts" bson:"ts"`
}

type CartItem struct {
	ID        string  `json:"id" bson:"id"`
	ProductID string  `json:"product_id" bson:"product_id"`
	VendorID  string  `json:"vendor_id" bson:"vendor_id"`
	Weight    float32 `json:"weight,omitempty" bson:"weight"`
	Quantity  int     `json:"quantity,omitempty" bson:"quantity"`
	TotalCost float64 `json:"total_cost" bson:"total_cost"`
}

func (c *CartItem) CalculateCartFee(fee *Fee) float64 {
	var totalCost float64

	if fee.ProductID == c.ProductID {
		if c.Weight != 0 {
			totalCost += float64(c.Weight) * fee.CostPerKg
		} else {
			totalCost += float64(c.Quantity) * fee.CostPerQty
		}
	} else {
		return 0
	}

	return totalCost
}

type CartStatuses string

const (
	CartActive   CartStatuses = "ACTIVE"   // cart has been created and active
	CartInactive CartStatuses = "INACTIVE" // cart has been inactivated and no longer active due to check out or session expiry
)

func IsValidCartStatus(status CartStatuses) bool {
	return status == CartActive || status == CartInactive
}

func SetCartStatus(status CartStatuses) (CartStatuses, error) {
	switch IsValidCartStatus(status) {
	case true:
		return status, nil
	default:
		return "", leetError.ErrorResponseBody(leetError.CartStatusesError, errors.New("invalid cart status"))
	}
}
