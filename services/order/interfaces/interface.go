package interfaces

import (
	"encoding/json"
	"github.com/leetatech/leeta_backend/services/order/domain"
	"net/http"
)

// CreateOrder godoc
// @Summary Create Order
// @Description The endpoint takes a domain.Order request and creates a new order
// @Tags Order
// @Accept json
// @Produce json
// @Param domain.Order body domain.Order true "create order request body"
// @Success 200 {object} domain.Order
// @Router /order/make_order [post]
func (handler *HTTPHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder domain.Order

	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		return
	}
	order, err := handler.OrderApplication.CreateOrder(r.Context(), newOrder)
	if err != nil {
		return
	}
	encodeResult(w, order, http.StatusOK)
}
