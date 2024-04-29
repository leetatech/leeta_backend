package interfaces

import (
	"encoding/json"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/helpers"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/query"
	"github.com/leetatech/leeta_backend/pkg/query/filter"
	"github.com/leetatech/leeta_backend/services/fees/application"
	"github.com/leetatech/leeta_backend/services/fees/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/samber/lo"
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
// @success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /fees/ [POST]
func (handler *FeesHttpHandler) CreateFeeHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.FeeQuotationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, leetError.ErrorResponseBody(leetError.UnmarshalError, err), http.StatusBadRequest)
		return
	}
	request, err = request.FeeTypeValidation()
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusBadRequest, err)
		return
	}

	response, err := handler.FeesApplication.FeeQuotation(r.Context(), request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	pkg.EncodeResult(w, response, http.StatusOK)
}

// FetchFeesHandler is the endpoint to all fees
// @Summary Get fees
// @Description The endpoint to get all types of fees
// @Tags Fees
// @Accept json
// @produce json
// @param query.ResultSelector body query.ResultSelector true "list fees request body"
// @Security BearerToken
// @success 200 {object} query.ResponseListWithMetadata[models.Fee]
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /fees/type [POST]
func (handler *FeesHttpHandler) FetchFeesHandler(w http.ResponseWriter, r *http.Request) {
	var request query.ResultSelector
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, leetError.ErrorResponseBody(leetError.UnmarshalError, err), http.StatusBadRequest)
		return
	}

	request, err = helpers.ValidateResultSelector(request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	fees, totalRecord, err := handler.FeesApplication.GetTypedFees(r.Context(), request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response := query.ResponseListWithMetadata[models.Fee]{
		Metadata: query.NewMetadata(request, totalRecord),
		Data:     fees,
	}
	pkg.EncodeResult(w, response, http.StatusOK)
}

// ListFeesOptions is the endpoint to get fees filter options
// @Summary Get fees filter options
// @Description Retrieve fees filter options
// @Tags Fees
// @Accept json
// @Produce json
// @Security BearerToken
// @Success 200 {object} filter.RequestOption
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /fees/options [get]
func (handler *FeesHttpHandler) ListFeesOptions(w http.ResponseWriter, r *http.Request) {
	requestOptions := lo.Map(listFeesOptions, toFilterOption)
	pkg.EncodeResult(w, requestOptions, http.StatusOK)
}

func toFilterOption(options filter.RequestOption, _ int) filter.RequestOption {
	return options
}
