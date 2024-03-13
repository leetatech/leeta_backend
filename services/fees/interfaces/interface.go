package interfaces

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/leetatech/leeta_backend/services/fees/application"
	"github.com/leetatech/leeta_backend/services/fees/domain"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/library/leetError"
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

// CreateFeeHandler is the endpoint to create fees
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
func (handler *FeesHttpHandler) CreateFeeHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.FeeQuotationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		library.EncodeResult(w, leetError.ErrorResponseBody(leetError.UnmarshalError, err), http.StatusBadRequest)
		return
	}

	response, err := handler.FeesApplication.FeeQuotation(r.Context(), request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	library.EncodeResult(w, response, http.StatusOK)
}

// GetFeesHandler is the endpoint to get fees
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
func (handler *FeesHttpHandler) GetFeesHandler(w http.ResponseWriter, r *http.Request) {
	response, err := handler.FeesApplication.GetFees(r.Context())
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	library.EncodeResult(w, response, http.StatusOK)
}

// GetFeeByProductIDHandler is the endpoint to get fees by product ID
// @Summary Get fee by product ID
// @Description The endpoint to get fees for gas refill by product ID
// @Tags fees
// @Accept json
// @produce json
// @param product_id path string true "product ID"
// @Security BearerToken
// @success 200 {object} library.DefaultResponse
// @Failure 401 {object} library.DefaultErrorResponse
// @Failure 400 {object} library.DefaultErrorResponse
// @Router /fees/product/{product_id} [GET]
func (handler *FeesHttpHandler) GetFeeByProductIDHandler(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "product_id")
	response, err := handler.FeesApplication.GetFeeByProductID(r.Context(), productID)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	library.EncodeResult(w, response, http.StatusOK)
}
