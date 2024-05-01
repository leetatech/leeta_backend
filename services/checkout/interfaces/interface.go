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

// UpdateGasRefillStatus is the endpoint used to update the status of a gas refill
// @Summary Update Gas refill request status
// @Description This endpoint is used to update the status of a gas refill (Cancel, Accept, Reject or Fulfill) request
// @Tags Gas Refill
// @Accept json
// @Produce json
// @Param domain.UpdateRefillRequest body domain.UpdateRefillRequest true "update gas refill by status request body"
// @Security BearerToken
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /gas-refill/ [put]
func (handler *CheckoutHttpHandler) UpdateGasRefillStatus(w http.ResponseWriter, r *http.Request) {

}

// GetGasRefill is the endpoint to get a single gas refill by id
// @Summary Gets a single gas refill
// @Description This is the endpoint to get the details of a single gas refill by refill-id
// @Tags Gas Refill
// @Accept json
// @produce json
// @Param			refill-id	path		string	true	"refill id"
// @Security BearerToken
// @success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /gas-refill/{refill_id} [get]
func (handler *CheckoutHttpHandler) GetGasRefill(w http.ResponseWriter, r *http.Request) {

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
