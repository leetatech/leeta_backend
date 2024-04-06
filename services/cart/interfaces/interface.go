package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/services/cart/application"
	"github.com/leetatech/leeta_backend/services/cart/domain"
	"net/http"
	"strconv"
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
// @Description The endpoint to delete items from cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerToken
// @Param cartID query string true "cartID"
// @Param cartItemID query string true "cartItemID"
// @Param productID query string true "productID"
// @Param weight query string false "weight"
// @Param quantity query string false "quantity"
// @Success 200 {object} library.DefaultResponse
// @Failure 401 {object} library.DefaultErrorResponse
// @Failure 400 {object} library.DefaultErrorResponse
// @Router /delete [delete]
func (handler *CartHttpHandler) DeleteCartItemHandler(w http.ResponseWriter, r *http.Request) {
	var (
		reducedWeightCount   float64
		reducedQuantityCount int
		err                  error
	)
	cartID := r.URL.Query().Get("cartID")
	cartItemID := r.URL.Query().Get("cartItemID")
	productID := r.URL.Query().Get("productID")
	weight := r.URL.Query().Get("weight")
	if weight != "" {
		reducedWeightCount, err = strconv.ParseFloat(weight, 32)
		if err != nil {
			library.EncodeResult(w, err, http.StatusBadRequest)
			return
		}
	}
	quantity := r.URL.Query().Get("quantity")
	if quantity != "" {
		reducedQuantityCount, _ = strconv.Atoi(quantity)
		if err != nil {
			library.EncodeResult(w, err, http.StatusBadRequest)
			return
		}
	}

	response, err := handler.CartApplication.DeleteCartItem(r.Context(), domain.DeleteCartItemRequest{
		CartID:               cartID,
		CartItemID:           cartItemID,
		ProductID:            productID,
		ReducedQuantityCount: reducedQuantityCount,
		ReducedWeightCount:   reducedWeightCount,
	})
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	library.EncodeResult(w, response, http.StatusOK)
}
