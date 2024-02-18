package interfaces

import (
	"github.com/leetatech/leeta_backend/services/gasrefill/application"
	_ "github.com/leetatech/leeta_backend/services/gasrefill/domain"
	_ "github.com/leetatech/leeta_backend/services/library"
	"net/http"
)

type GasRefillHttpHandler struct {
	GasRefillApplication application.GasRefillApplication
}

func NewGasRefillHTTPHandler(refillApplication application.GasRefillApplication) *GasRefillHttpHandler {
	return &GasRefillHttpHandler{
		GasRefillApplication: refillApplication,
	}
}

// RequestRefill is the endpoint to handle gas refill
// @Summary Request gas refill
// @Description The endpoint to request for a gas refill
// @Tags GasRefill
// @Accept json
// @Produce json
// @Param domain.GasRefillRequest body domain.GasRefillRequest true "Gas refill request body"
// @Security BearerToken
// @Success 200 {object} domain.Gas
// @Failure 401 {object} library.DefaultErrorResponse
// @Failure 400 {object} library.DefaultErrorResponse
// @Router /gas-refill [post]
func (handler *GasRefillHttpHandler) RequestRefill(w http.ResponseWriter, r *http.Request) {

}

// UpdateGasRefillStatus is the endpoint used to update the status of a gas refill
// @Summary Update Gas refill request status
// @Description This endpoint is used to update the status of a gas refill (Cancel, Accept, Reject or Fulfill) request
// @Tags GasRefill
// @Accept json
// @Produce json
// @Param domain.UpdateRefillRequest body domain.UpdateRefillRequest true "update gas refill by status request body"
// @Security BearerToken
// @Success 200 {object} domain.Gas
// @Failure 401 {object} library.DefaultErrorResponse
// @Failure 400 {object} library.DefaultErrorResponse
// @Router /gas-refill/ [put]
func (handler *GasRefillHttpHandler) UpdateGasRefillStatus(w http.ResponseWriter, r *http.Request) {

}

// GetGasRefill is the endpoint to get a single gas refill by id
// @Summary Gets a single gas refill
// @Description This is the endpoint to get the details of a single gas refill by refill-id
// @Tags GasRefill
// @Accept json
// @produce json
// @Param			refill-id	path		string	true	"refill id"
// @Security BearerToken
// @success 200 {object} domain.Gas
// @Failure 401 {object} library.DefaultErrorResponse
// @Failure 400 {object} library.DefaultErrorResponse
// @Router /gas-refill/{refill_id} [get]
func (handler *GasRefillHttpHandler) GetGasRefill(w http.ResponseWriter, r *http.Request) {

}

// ListRefill handles listing all refill requests
// @Summary List all gas refill requests
// @Description The endpoint takes the order status, pages and limit and then returns the requested orders
// @Tags GasRefill
// @Accept json
// @produce json
// @param domain.ListRefillFilter body domain.ListRefillFilter true "get refill by status, use filter for filtering responses (not implemented)"
// @Security BearerToken
// @success 200 {object} []domain.Gas
// @Failure 401 {object} library.DefaultErrorResponse
// @Failure 400 {object} library.DefaultErrorResponse
// @Router /gas-refill/list [POST]
func (handler *GasRefillHttpHandler) ListRefill(w http.ResponseWriter, r *http.Request) {
}
