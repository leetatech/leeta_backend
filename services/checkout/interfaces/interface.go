package interfaces

import (
	"encoding/json"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/services/checkout/application"
	"github.com/leetatech/leeta_backend/services/checkout/domain"
	"net/http"
)

type CheckoutHttpHandler struct {
	CheckoutApplication application.CheckoutApplication
}

func NewCheckoutHTTPHandler(checkoutApplication application.CheckoutApplication) *CheckoutHttpHandler {
	return &CheckoutHttpHandler{
		CheckoutApplication: checkoutApplication,
	}
}

// Checkout is the endpoint to check out from cart
// @Summary Check out from cart
// @Description The endpoint to allows the user to check out from the cart
// @Tags Checkout
// @Accept json
// @Produce json
// @Param domain.CheckoutRequest body domain.CheckoutRequest true "Check out request body"
// @Security BearerToken
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /checkout/ [post]
func (handler *CheckoutHttpHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var request domain.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response, err := handler.CheckoutApplication.Checkout(r.Context(), request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	pkg.EncodeResult(w, response, http.StatusOK)
}

// UpdateCheckout is the endpoint used to update the status of the checkout
// @Summary Update checkout status
// @Description This endpoint is used to update the status of a checkout (Cancel, Accept, Reject or Fulfill) by customers, vendor and admin
// @Tags Checkout
// @Accept json
// @Produce json
// @Param domain.UpdateCheckoutRequest body domain.UpdateCheckoutRequest true "update checkout request body"
// @Security BearerToken
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /checkout/ [put]
func (handler *CheckoutHttpHandler) UpdateCheckout(w http.ResponseWriter, r *http.Request) {
	var request domain.UpdateCheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response, err := handler.CheckoutApplication.UpdateCheckout(r.Context(), request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	pkg.EncodeResult(w, response, http.StatusOK)
}

// ListRefill handles listing all refill requests
// @Summary List all gas refill requests
// @Description The endpoint takes the order status, pages and limit and then returns the requested orders
// @Tags Gas Refill
// @Accept json
// @produce json
// @param domain.ListRefillFilter body domain.ListRefillFilter true "get refill by status, use filter for filtering responses (not implemented)"
// @Security BearerToken
// @success 200 {object} []pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /gas-refill/list [POST]
func (handler *CheckoutHttpHandler) ListRefill(w http.ResponseWriter, r *http.Request) {
}
