package interfaces

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	_ "github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/pkg/helpers"
	"github.com/leetatech/leeta_backend/pkg/jwtmiddleware"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/order/application"
	"github.com/leetatech/leeta_backend/services/order/domain"
	"github.com/leetatech/leeta_backend/services/web"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"net/http"
)

type OrderHttpHandler struct {
	OrderApplication application.Order
}

func New(orderApplication application.Order) *OrderHttpHandler {
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
// @Param domain.UpdateStatusRequest body domain.UpdateStatusRequest true "update order by status request body"
// @Security BearerToken
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /order/status [put]
func (handler *OrderHttpHandler) UpdateOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.UpdateStatusRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	resp, err := handler.OrderApplication.UpdateOrderStatus(r.Context(), request)
	if err != nil {
		log.Err(err).Msg("error updating order status")
		helpers.CheckErrorType(err, w)
		return
	}

	jwtmiddleware.WriteJSONResponse(w, resp, http.StatusOK)
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

	jwtmiddleware.WriteJSONResponse(w, order, http.StatusOK)
}

// GetCustomerOrdersByStatusHandler godoc
// @Summary Get Customer Order By Status
// @Description The endpoint takes the order status, pages and limit and then returns the requested orders
// @Tags Order
// @Accept json
// @produce json
// @param domain.GetCustomerOrdersRequest body domain.GetCustomerOrdersRequest true "get customer orders by status request body"
// @Security BearerToken
// @success 200 {object} []domain.Response
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /order/ [get]
func (handler *OrderHttpHandler) GetCustomerOrdersByStatusHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.GetCustomerOrdersRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		jwtmiddleware.WriteJSONResponse(w, err, http.StatusBadRequest)
		return
	}
	orders, err := handler.OrderApplication.GetCustomerOrdersByStatus(r.Context(), request)
	if err != nil {
		helpers.CheckErrorType(err, w)
		return
	}
	jwtmiddleware.WriteJSONResponse(w, orders, http.StatusOK)
}

// ListOrdersHandler is the endpoint to list all orders
// @Summary List orders
// @Description The endpoint to list all orders. List endpoint can be configured with the filters
// @Tags Order
// @Accept json
// @produce json
// @param query.ResultSelector body query.ResultSelector true "list orders request body"
// @Security BearerToken
// @success 200 {object} query.ResponseListWithMetadata[models.Order]
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /order/ [PUT]
func (handler *OrderHttpHandler) ListOrdersHandler(w http.ResponseWriter, r *http.Request) {
	resultSelector, err := web.PrepareResultSelector(r, listOrdersOptions, allowedSortFields, web.ResultSelectorDefaults(defaultSortingRequest))
	if err != nil {
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusBadRequest, errs.Body(errs.InvalidRequestError, err))
		return
	}

	orders, totalRecord, err := handler.OrderApplication.ListOrders(r.Context(), resultSelector)
	if err != nil {
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusInternalServerError, errs.Body(errs.InternalError, err))
		return
	}

	response := query.ResponseListWithMetadata[models.Order]{
		Metadata: query.NewMetadata(resultSelector, totalRecord),
		Data:     orders,
	}
	jwtmiddleware.WriteJSONResponse(w, response, http.StatusOK)
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
	jwtmiddleware.WriteJSONResponse(w, requestOptions, http.StatusOK)
}

// ListOrderStatusHistoryHandler godoc
// @Summary Get Order Status History
// @Description The endpoint takes the order id and then returns the requested order status history
// @Tags Order
// @Accept json
// @produce json
// @Param			order_id	path		string	true	"order id"
// @Security BearerToken
// @success 200 {object} []models.StatusHistory
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /order/status/history/{order_id} [get]
func (handler *OrderHttpHandler) ListOrderStatusHistoryHandler(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")
	statusHistory, err := handler.OrderApplication.ListOrderStatusHistory(r.Context(), orderID)
	if err != nil {
		helpers.CheckErrorType(err, w)
		return
	}

	jwtmiddleware.WriteJSONResponse(w, statusHistory, http.StatusOK)
}

func toFilterOption(options filter.RequestOption, _ int) filter.RequestOption {
	return options
}
