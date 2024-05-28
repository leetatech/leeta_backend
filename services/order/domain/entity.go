package domain

import "github.com/leetatech/leeta_backend/services/models"

type OrderRequest struct {
	ProductID string `json:"product_id" bson:"product_id"`
} // @name OrderRequest

type OrderResponse struct {
	ID              string               `json:"id" bson:"id"`
	ProductID       string               `json:"product_id" bson:"product_id"`
	CustomerID      string               `json:"customer_id" bson:"customer_id"`
	VendorID        string               `json:"vendor_id" bson:"vendor_id"`
	VAT             float64              `json:"vat" bson:"vat"`
	DeliveryFee     float64              `json:"delivery_fee" bson:"delivery_fee"`
	Total           float64              `json:"total" bson:"total"`
	Status          models.OrderStatuses `json:"status" bson:"status"`
	StatusTs        int64                `json:"status_ts" bson:"status_ts"`
	Ts              int64                `json:"ts" bson:"ts"`
	CustomerDetails models.Admin         `json:"customer" bson:"customer"`
	ProductDetails  models.Product       `json:"product_details"`
} // @name OrderResponse

type UpdateOrderStatusRequest struct {
	OrderId     string               `json:"order_id" bson:"order_id"`
	OrderStatus models.OrderStatuses `json:"order_status" bson:"order_status"`
	Reason      string               `json:"reason" bson:"reason"`
} // @name UpdateOrderStatusRequest

type PersistOrderUpdate struct {
	UpdateOrderStatusRequest
	StatusHistory models.StatusHistory `json:"status_history" bson:"status_history"`
}

type GetCustomerOrders struct {
	UserId string `json:"user_id" bson:"user_id"`
	GetCustomerOrdersRequest
} // @name GetCustomerOrders

type GetCustomerOrdersRequest struct {
	OrderStatus []models.OrderStatuses `json:"order_status" bson:"order_status"`
	Limit       int64                  `json:"limit" bson:"limit"`
	Page        int64                  `json:"page" bson:"page"`
} // @name GetCustomerOrdersRequest
