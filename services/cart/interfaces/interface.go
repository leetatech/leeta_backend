package interfaces

import (
	"encoding/json"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/services/cart/application"
	"github.com/leetatech/leeta_backend/services/cart/domain"
	"github.com/rs/zerolog/log"
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

// AddToCartHandler is the endpoint to add items to cart
// @Summary Add items to cart
// @Description The endpoint to add items to cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param domain.CartItem body domain.CartItem true "add to cart request body"
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /cart/add [post]
func (handler *CartHttpHandler) AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.CartItem
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response, err := handler.CartApplication.AddToCart(r.Context(), request)
	if err != nil {
		log.Debug().Msgf("error adding item to cart %v", err)
		pkg.EncodeResult(w, err, http.StatusInternalServerError)
		return
	}
	pkg.EncodeResult(w, response, http.StatusOK)
}

// InactivateCartHandler is the endpoint to inactivate carts
// @Summary Request cart inactivation
// @Description The endpoint to request for a cart inactivation
// @Tags Cart
// @Accept json
// @Produce json
// @Param domain.InactivateCart body domain.InactivateCart true "inactivate cart request body"
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /cart/inactivate [put]
func (handler *CartHttpHandler) InactivateCartHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.InactivateCart
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response, err := handler.CartApplication.InactivateCart(r.Context(), request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	pkg.EncodeResult(w, response, http.StatusOK)
}

// UpdateCartItemQuantityHandler is the endpoint to increase cart item quantity
// @Summary increase or reduce cart item quantity
// @Description The endpoint to increase or reduce cart item quantity
// @Tags Cart
// @Accept json
// @Produce json
// @Param domain.UpdateCartItemQuantityRequest body domain.UpdateCartItemQuantityRequest true "update cart item quantity request body"
// @Security BearerToken
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /cart/item [put]
func (handler *CartHttpHandler) UpdateCartItemQuantityHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.UpdateCartItemQuantityRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	response, err := handler.CartApplication.UpdateCartItemQuantity(r.Context(), request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	pkg.EncodeResult(w, response, http.StatusOK)
}
