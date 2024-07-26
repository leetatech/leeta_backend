package interfaces

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	_ "github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/pkg/jwtmiddleware"
	"github.com/leetatech/leeta_backend/services/cart/application"
	"github.com/leetatech/leeta_backend/services/cart/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/web"
	"github.com/samber/lo"
	"net/http"
)

type CartHttpHandler struct {
	CartApplication application.Cart
}

func New(cartApplication application.Cart) *CartHttpHandler {
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
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	response, err := handler.CartApplication.Add(r.Context(), request)
	if err != nil {
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	jwtmiddleware.WriteJSONResponse(w, response, http.StatusCreated)
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
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusBadRequest, errors.New("cart_id is required"))
		return
	}

	err := handler.CartApplication.Delete(r.Context(), cartID)
	if err != nil {
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("error deleting cart: %w", err))
		return
	}
	jwtmiddleware.WriteJSONResponse(w, nil, http.StatusAccepted)
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
		jwtmiddleware.WriteJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	if isValid, err := request.IsValid(); !isValid {
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusBadRequest, fmt.Errorf("requst is not valid: %w", err))
		return
	}

	response, err := handler.CartApplication.UpdateItemQuantity(r.Context(), request)
	if err != nil {
		jwtmiddleware.WriteJSONResponse(w, err, http.StatusInternalServerError)
		return
	}

	jwtmiddleware.WriteJSONResponse(w, response, http.StatusAccepted)
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
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusBadRequest, fmt.Errorf("cart_item_id is required"))
		return
	}

	err := handler.CartApplication.DeleteItem(r.Context(), cartItemID)
	if err != nil {
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	jwtmiddleware.WriteJSONResponse(w, nil, http.StatusAccepted)
}

// ListCart is the endpoint to list cart.
// @Summary List cart and items. Use result selector to filter results and manage pagination
// @Description The endpoint to get a user cart, and the items in the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerToken
// @Param query.ResultSelector body query.ResultSelector true "list cart request body"
// @Success 200 {object} query.ResponseWithMetadata[models.Cart]
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /cart [put]
func (handler *CartHttpHandler) ListCart(w http.ResponseWriter, r *http.Request) {
	resultSelector, err := web.PrepareResultSelector(r, listCartOptions, allowedSortFields, web.ResultSelectorDefaults(defaultSortingRequest))
	if err != nil {
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusBadRequest, errs.Body(errs.InvalidRequestError, err))
		return
	}

	carts, totalRecord, err := handler.CartApplication.ListCart(r.Context(), resultSelector)
	if err != nil {
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	response := query.ResponseWithMetadata[models.Cart]{
		Metadata: query.NewMetadata(resultSelector, totalRecord),
		Data:     carts,
	}

	jwtmiddleware.WriteJSONResponse(w, response, http.StatusOK)
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
	jwtmiddleware.WriteJSONResponse(w, requestOptions, http.StatusOK)
}

func toFilterOption(options filter.RequestOption, _ int) filter.RequestOption {
	return options
}

// Checkout is the endpoint to check out from cart
// @Summary Check out from cart
// @Description The endpoint to allows the user to check out from the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param domain.CartCheckoutRequest body domain.CartCheckoutRequest true "Cart checkout request body"
// @Security BearerToken
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /cart/checkout [post]
func (handler *CartHttpHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var request domain.CartCheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		jwtmiddleware.WriteJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	response, err := handler.CartApplication.Checkout(r.Context(), request)
	if err != nil {
		jwtmiddleware.WriteJSONResponse(w, err, http.StatusBadRequest)
		return
	}
	jwtmiddleware.WriteJSONResponse(w, response, http.StatusOK)
}
