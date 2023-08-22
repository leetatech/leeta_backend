package domain

import "github.com/leetatech/leeta_backend/services/library/models"

type OrderRequest struct {
	ID        string `json:"id" bson:"id"`
	ProductID string `json:"product_id" bson:"product_id"`
	VendorID  string `json:"vendor_id" bson:"vendor_id"`
	Status    string `json:"status" bson:"status"`
	StatusTs  int64  `json:"status_ts" bson:"status_ts"`
	Ts        int64  `json:"ts" bson:"ts"`
} // @name OrderRequest

type OrderResponse struct {
	OrderDetails    models.Order    `json:"order_details"`
	CustomerDetails models.Customer `json:"customer_details"`
	ProductDetails  models.Product  `json:"product_details"`
} // @name OrderResponse
