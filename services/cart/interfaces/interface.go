package interfaces

import (
	"encoding/json"
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
// @Security BearerToken
// @Param domain.AddToCartRequest body domain.AddToCartRequest true "add to cart request body"
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /cart/add [post]
func (handler *CartHttpHandler) AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.AddToCartRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response, err := handler.CartApplication.AddToCart(r.Context(), request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	pkg.EncodeResult(w, response, http.StatusOK)
}

// DeleteCartHandler is the endpoint to delete carts
// @Summary Request cart deletion
// @Description The endpoint to request for a cart deletion
// @Tags Cart
// @Accept json
// @Produce json
// @Param cartID query string true "cartID"
// @Security BearerToken
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /cart/ [delete]
func (handler *CartHttpHandler) DeleteCartHandler(w http.ResponseWriter, r *http.Request) {
	cartID := r.URL.Query().Get("cartID")

	response, err := handler.CartApplication.DeleteCart(r.Context(), domain.DeleteCartRequest{
		ID: cartID,
	})
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
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
// @Param cartItemID query string true "cartItemID"
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /cart/item [delete]
func (handler *CartHttpHandler) DeleteCartItemHandler(w http.ResponseWriter, r *http.Request) {
	cartItemID := r.URL.Query().Get("cartItemID")

	response, err := handler.CartApplication.DeleteCartItem(r.Context(), domain.DeleteCartItemRequest{CartItemID: cartItemID})
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	pkg.EncodeResult(w, response, http.StatusOK)
}
