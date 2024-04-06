package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/services/cart/application"
	"github.com/leetatech/leeta_backend/services/cart/domain"
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
		pkg.EncodeErrorResult(w, http.StatusBadRequest, err)
		return
	}

	response, err := handler.CartApplication.AddToCart(r.Context(), request)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
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
// @deprecated
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
// @Router /cart/item/quantity [put]
func (handler *CartHttpHandler) UpdateCartItemQuantityHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.UpdateCartItemQuantityRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	if isValid, err := request.IsValid(); !isValid {
		pkg.EncodeErrorResult(w, http.StatusBadRequest, fmt.Errorf("requst is not valid: %w", err))
		return
	}

	response, err := handler.CartApplication.UpdateCartItemQuantity(r.Context(), request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusInternalServerError)
		return
	}

	pkg.EncodeResult(w, response, http.StatusOK)
}

// DeleteCartItemHandler is the endpoint to delete items from cart
// @Summary Delete items from cart
// @Description The endpoint to delete items from cart. This endpoint also deletes an entire cart if there is no item left in the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerToken
// @Param cart_item_id query string true "cart_item_id"
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /cart/item/{cart_item_id} [delete]
func (handler *CartHttpHandler) DeleteCartItemHandler(w http.ResponseWriter, r *http.Request) {
	cartItemID := chi.URLParam(r, "cart_item_id")
	if cartItemID == "" {
		pkg.EncodeErrorResult(w, http.StatusBadRequest, fmt.Errorf("cart_item_id is required"))
		return
	}

	response, err := handler.CartApplication.DeleteCartItem(r.Context(), cartItemID)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
		return
	}

	pkg.EncodeResult(w, response, http.StatusOK)
}
