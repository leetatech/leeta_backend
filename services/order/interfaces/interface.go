package interfaces

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/helpers"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/query"
	"github.com/leetatech/leeta_backend/pkg/query/filter"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/order/application"
	"github.com/leetatech/leeta_backend/services/order/domain"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
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

// UpdateOrderStatusHandler godoc
// @Summary Update Order Status
// @Description The endpoint takes the order update request and updates the status of the order
// @Tags Order
// @Accept json
// @Produce json
// @Param domain.UpdateOrderStatusRequest body domain.UpdateOrderStatusRequest true "update order by status request body"
// @Security BearerToken
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /order/status [put]
func (handler *OrderHttpHandler) UpdateOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.UpdateOrderStatusRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusBadRequest, err)
		return
	}

	resp, err := handler.OrderApplication.UpdateOrderStatus(r.Context(), request)
	if err != nil {
		log.Err(err).Msg("error updating order status")
		helpers.CheckErrorType(err, w)
		return
	}

	pkg.EncodeResult(w, resp, http.StatusOK)
}

// GetOrderByIDHandler godoc
// @Summary Get Customer Order By id
// @Description The endpoint takes the order id and then returns the requested order
// @Tags Order
// @Accept json
// @produce json
// @Param			order_id	path		string	true	"order id"
// @Security BearerToken
// @success 200 {object} models.Order
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /order/id/{order_id} [get]
func (handler *OrderHttpHandler) GetOrderByIDHandler(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")
	var order *models.Order
	order, err := handler.OrderApplication.GetOrderByID(r.Context(), orderID)
	if err != nil {
		helpers.CheckErrorType(err, w)
		return
	}

	pkg.EncodeResult(w, order, http.StatusOK)
}

// GetCustomerOrdersByStatusHandler godoc
// @Summary Get Customer Order By Status
// @Description The endpoint takes the order status, pages and limit and then returns the requested orders
// @Tags Order
// @Accept json
// @produce json
// @param domain.GetCustomerOrdersRequest body domain.GetCustomerOrdersRequest true "get customer orders by status request body"
// @Security BearerToken
// @success 200 {object} []domain.OrderResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /order/ [get]
func (handler *OrderHttpHandler) GetCustomerOrdersByStatusHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.GetCustomerOrdersRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	orders, err := handler.OrderApplication.GetCustomerOrdersByStatus(r.Context(), request)
	if err != nil {
		helpers.CheckErrorType(err, w)
		return
	}
	pkg.EncodeResult(w, orders, http.StatusOK)
}

// FetchOrdersHandler is the endpoint to all orders
// @Summary Get orders
// @Description The endpoint to get all orders using several filters
// @Tags Order
// @Accept json
// @produce json
// @param query.ResultSelector body query.ResultSelector true "list orders request body"
// @Security BearerToken
// @success 200 {object} query.ResponseListWithMetadata[models.Order]
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /order/ [POST]
func (handler *OrderHttpHandler) FetchOrdersHandler(w http.ResponseWriter, r *http.Request) {
	var request query.ResultSelector
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, leetError.ErrorResponseBody(leetError.UnmarshalError, err), http.StatusBadRequest)
		return
	}

	request, err = helpers.ValidateResultSelector(request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	orders, totalRecord, err := handler.OrderApplication.ListOrders(r.Context(), request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response := query.ResponseListWithMetadata[models.Order]{
		Metadata: query.NewMetadata(request, totalRecord),
		Data:     orders,
	}
	pkg.EncodeResult(w, response, http.StatusOK)
}

// ListOrdersOptions is the endpoint to get orders filter options
// @Summary Get orders filter options
// @Description Retrieve orders filter options
// @Tags Order
// @Accept json
// @Produce json
// @Security BearerToken
// @Success 200 {object} filter.RequestOption
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /order/options [get]
func (handler *OrderHttpHandler) ListOrdersOptions(w http.ResponseWriter, r *http.Request) {
	requestOptions := lo.Map(listOrdersOptions, toFilterOption)
	pkg.EncodeResult(w, requestOptions, http.StatusOK)
}

func toFilterOption(options filter.RequestOption, _ int) filter.RequestOption {
	return options
}
