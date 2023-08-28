package domain

import "github.com/leetatech/leeta_backend/services/library/models"

type OrderRequest struct {
	ID        string `json:"id" bson:"id"`
	ProductID string `json:"product_id" bson:"product_id"`
	Status    string `json:"status" bson:"status"`
	StatusTs  int64  `json:"status_ts" bson:"status_ts"`
	Ts        int64  `json:"ts" bson:"ts"`
} // @name OrderRequest

type OrderResponse struct {
	OrderDetails    models.Order    `json:"order" bson:"order"`
	CustomerDetails models.Customer `json:"customer" bson:"customers"`
	ProductDetails  models.Product  `json:"product_details"`
} // @name OrderResponse

type UpdateOrderStatusRequest struct {
	OrderId     string               `json:"order_id" bson:"order_id"`
	OrderStatus models.OrderStatuses `json:"order_status" bson:"order_status"`
} // @name UpdateOrderStatusRequest

type GetCustomerOrders struct {
	UserId string `json:"user_id" bson:"user_id"`
	GetCustomerOrdersRequest
} // @name GetCustomerOrders

type GetCustomerOrdersRequest struct {
	OrderStatus []models.OrderStatuses `json:"order_status" bson:"order_status"`
	Limit       int64                  `json:"limit" bson:"limit"`
	Page        int64                  `json:"page" bson:"page"`
} // @name GetCustomerOrdersRequest
