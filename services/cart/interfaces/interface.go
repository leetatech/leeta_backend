package interfaces

import (
	"encoding/json"
	"github.com/leetatech/leeta_backend/services/cart/application"
	"github.com/leetatech/leeta_backend/services/cart/domain"
	"github.com/leetatech/leeta_backend/services/library"
	"net/http"
)

type CartHttpHandler struct {
	CartApplication application.CartApplication
}

func NewCartHTTPHandler(cartApplication application.CartApplication) *CartHttpHandler {
	return &CartHttpHandler{
		CartApplication: cartApplication,
	}
}

// InactivateCartHandler is the endpoint to inactivate carts
// @Summary Request cart inactivation
// @Description The endpoint to request for a cart inactivation
// @Tags Cart
// @Accept json
// @Produce json
// @Param domain.InactivateCart body domain.InactivateCart true "inactivate cart request body"
// @Success 200 {object} library.DefaultResponse
// @Failure 401 {object} library.DefaultErrorResponse
// @Failure 400 {object} library.DefaultErrorResponse
// @Router /inactivate [put]
func (handler *CartHttpHandler) InactivateCartHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.InactivateCart
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response, err := handler.CartApplication.InactivateCart(r.Context(), request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	library.EncodeResult(w, response, http.StatusOK)
}
