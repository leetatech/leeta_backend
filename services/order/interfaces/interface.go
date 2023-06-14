package interfaces

import (
	"encoding/json"
	"github.com/leetatech/leeta_backend/services/order/domain"
	"net/http"
)

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
