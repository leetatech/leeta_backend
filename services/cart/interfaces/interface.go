package interfaces

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/helpers"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/query"
	"github.com/leetatech/leeta_backend/pkg/query/filter"
	"github.com/leetatech/leeta_backend/services/cart/application"
	"github.com/leetatech/leeta_backend/services/cart/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/samber/lo"
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

// AddToCart is the endpoint to add items to cart
// @Summary Add items to cart
// @Description The endpoint to add items to cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerToken
// @Param domain.CartItem body domain.CartItem true "add to cart request body"
// @Success 201 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /cart/add [post]
func (handler *CartHttpHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
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
	pkg.EncodeResult(w, response, http.StatusCreated)
}

// DeleteCart is the endpoint to delete carts
// @Summary Delete item from a cart
// @Description The endpoint is used to delete an item from a cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param cartID query string true "cartID"
// @Security BearerToken
// @Success 202 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /cart/{cart_id} [delete]
func (handler *CartHttpHandler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	cartID := chi.URLParam(r, "cart_id")
	if cartID == "" {
		pkg.EncodeErrorResult(w, http.StatusBadRequest, errors.New("cart_id is required"))
		return
	}

	err := handler.CartApplication.DeleteCart(r.Context(), cartID)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, fmt.Errorf("error deleting cart: %w", err))
		return
	}
	pkg.EncodeResult(w, nil, http.StatusAccepted)
}

// UpdateCartItemQuantity is the endpoint to increase cart item quantity
// @Summary increase or reduce cart item quantity
// @Description The endpoint to increase or reduce cart item quantity
// @Tags Cart
// @Accept json
// @Produce json
// @Param domain.UpdateCartItemQuantityRequest body domain.UpdateCartItemQuantityRequest true "update cart item quantity request body"
// @Security BearerToken
// @Success 202 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /cart/item/quantity [put]
func (handler *CartHttpHandler) UpdateCartItemQuantity(w http.ResponseWriter, r *http.Request) {
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

	pkg.EncodeResult(w, response, http.StatusAccepted)
}

// DeleteCartItem is the endpoint to delete items from cart
// @Summary Delete items from cart
// @Description The endpoint to delete items from cart. This endpoint also deletes an entire cart if there is no item left in the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerToken
// @Param cart_item_id query string true "cart_item_id"
// @Success 202 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /cart/item/{cart_item_id} [delete]
func (handler *CartHttpHandler) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	cartItemID := chi.URLParam(r, "cart_item_id")
	if cartItemID == "" {
		pkg.EncodeErrorResult(w, http.StatusBadRequest, fmt.Errorf("cart_item_id is required"))
		return
	}

	err := handler.CartApplication.DeleteCartItem(r.Context(), cartItemID)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
		return
	}

	pkg.EncodeResult(w, nil, http.StatusAccepted)
}

// ListCart is the endpoint to list cart.
// @Summary Get a user cart and list the cart items. Use result selector to filter results and manage pagination
// @Description The endpoint to get a user cart, and the items in the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerToken
// @Param query.ResultSelector body query.ResultSelector true "list cart request body"
// @Success 200 {object} query.ResponseWithMetadata[models.Cart]
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /cart [post]
func (handler *CartHttpHandler) ListCart(w http.ResponseWriter, r *http.Request) {
	var request query.ResultSelector

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusBadRequest, leetError.ErrorResponseBody(leetError.UnmarshalError, err))
		return
	}

	request, err = helpers.ValidateResultSelector(request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	carts, totalRecord, err := handler.CartApplication.ListCart(r.Context(), request)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
		return
	}

	response := query.ResponseWithMetadata[models.Cart]{
		Metadata: query.NewMetadata(request, totalRecord),
		Data:     carts,
	}

	pkg.EncodeResult(w, response, http.StatusOK)
}

// ListCartOptions is the endpoint to get cart filter options
// @Summary Get cart filter options
// @Description Retrieve cart filter options
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerToken
// @Success 200 {object} filter.RequestOption
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /cart/options [get]
func (handler *CartHttpHandler) ListCartOptions(w http.ResponseWriter, r *http.Request) {
	requestOptions := lo.Map(listCartOptions, toFilterOption)
	pkg.EncodeResult(w, requestOptions, http.StatusOK)
}

func toFilterOption(options filter.RequestOption, _ int) filter.RequestOption {
	return options
}
