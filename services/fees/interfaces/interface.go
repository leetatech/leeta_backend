package interfaces

import (
	"encoding/json"
	"github.com/leetatech/leeta_backend/services/fees/application"
	"github.com/leetatech/leeta_backend/services/fees/domain"
	"github.com/leetatech/leeta_backend/services/library"
	"net/http"
)

type FeesHttpHandler struct {
	FeesApplication application.FeesApplication
}

func NewFeesHTTPHandler(feesApplication application.FeesApplication) *FeesHttpHandler {
	return &FeesHttpHandler{
		FeesApplication: feesApplication,
	}
}

// CreateFees is the endpoint to create fees
// @Summary Create fees
// @Description The endpoint to create fees for gas refill
// @Tags Fees
// @Accept json
// @produce json
// @param domain.FeeQuotationRequest body domain.FeeQuotationRequest true "create fees request body"
// @Security BearerToken
// @success 200 {object} library.DefaultResponse
// @Failure 401 {object} library.DefaultErrorResponse
// @Failure 400 {object} library.DefaultErrorResponse
// @Router /fees/ [POST]
func (handler *FeesHttpHandler) CreateFees(w http.ResponseWriter, r *http.Request) {
	var request domain.FeeQuotationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response, err := handler.FeesApplication.FeeQuotation(r.Context(), request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	library.EncodeResult(w, response, http.StatusOK)
}

// GetFees is the endpoint to get fees
// @Summary Get fees
// @Description The endpoint to get fees for gas refill
// @Tags Fees
// @Accept json
// @produce json
// @Security BearerToken
// @success 200 {object} library.DefaultResponse
// @Failure 401 {object} library.DefaultErrorResponse
// @Failure 400 {object} library.DefaultErrorResponse
// @Router /fees/ [GET]
func (handler *FeesHttpHandler) GetFees(w http.ResponseWriter, r *http.Request) {
	response, err := handler.FeesApplication.GetFees(r.Context())
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	library.EncodeResult(w, response, http.StatusOK)
}
