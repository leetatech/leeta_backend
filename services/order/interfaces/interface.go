package interfaces

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/order/application"
	"github.com/leetatech/leeta_backend/services/order/domain"
	"net/http"
)

type OrderHttpHandler struct {
	OrderApplication application.OrderApplication
}

func NewOrderHTTPHandler(orderApplication application.OrderApplication) *OrderHttpHandler {
	return &OrderHttpHandler{
		OrderApplication: orderApplication,
	}

}

// CreateOrderHandler godoc
// @Summary Create Order
// @Description The endpoint takes the order request and creates a new order
// @Tags Order
// @Accept json
// @Produce json
// @Param domain.OrderRequest body domain.OrderRequest true "create order request body"
// @Security BearerToken
// @Success 200 {object} domain.OrderRequest
// @Failure 401 {object} library.DefaultErrorResponse
// @Failure 400 {object} library.DefaultErrorResponse
// @Router /order/make_order [post]
func (handler *OrderHttpHandler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var newOrder domain.OrderRequest

	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	resp, err := handler.OrderApplication.CreateOrder(r.Context(), newOrder)
	if err != nil {
		library.CheckErrorType(err, w)
		return
	}
	library.EncodeResult(w, resp, http.StatusOK)
}

// UpdateOrderStatusHandler godoc
// @Summary Update Order Status
// @Description The endpoint takes the order update request and updates the status of the order
// @Tags Order
// @Accept json
// @Produce json
// @Param domain.UpdateOrderStatusRequest body domain.UpdateOrderStatusRequest true "update order by status request body"
// @Security BearerToken
// @Success 200 {object} library.DefaultResponse
// @Failure 401 {object} library.DefaultErrorResponse
// @Failure 400 {object} library.DefaultErrorResponse
// @Router /order/status [put]
func (handler *OrderHttpHandler) UpdateOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.UpdateOrderStatusRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	resp, err := handler.OrderApplication.UpdateOrderStatus(r.Context(), request)
	if err != nil {
		library.CheckErrorType(err, w)
		return
	}

	library.EncodeResult(w, resp, http.StatusOK)
}

// GetOrderByIDHandler godoc
// @Summary Get Customer Order By id
// @Description The endpoint takes the order id and then returns the requested order
// @Tags Order
// @Accept json
// @produce json
// @Param			order_id	path		string	true	"order id"
// @Security BearerToken
// @success 200 {object} []domain.OrderResponse
// @Failure 401 {object} library.DefaultErrorResponse
// @Failure 400 {object} library.DefaultErrorResponse
// @Router /order/id/{order_id} [get]
func (handler *OrderHttpHandler) GetOrderByIDHandler(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")

	order, err := handler.OrderApplication.GetOrderByID(r.Context(), orderID)
	if err != nil {
		library.CheckErrorType(err, w)
		return
	}

	library.EncodeResult(w, order, http.StatusOK)
}

// GetCustomerOrdersByStatusHandler godoc
// @Summary Get Customer Orders By Status
// @Description The endpoint takes the order status, pages and limit and then returns the requested orders
// @Tags Order
// @Accept json
// @produce json
// @param domain.GetCustomerOrdersRequest body domain.GetCustomerOrdersRequest true "get customer orders by status request body"
// @Security BearerToken
// @success 200 {object} []domain.OrderResponse
// @Failure 401 {object} library.DefaultErrorResponse
// @Failure 400 {object} library.DefaultErrorResponse
// @Router /order/ [get]
func (handler *OrderHttpHandler) GetCustomerOrdersByStatusHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.GetCustomerOrdersRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	orders, err := handler.OrderApplication.GetCustomerOrdersByStatus(r.Context(), request)
	if err != nil {
		library.CheckErrorType(err, w)
		return
	}
	library.EncodeResult(w, orders, http.StatusOK)
}
