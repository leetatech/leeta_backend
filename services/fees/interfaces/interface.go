package interfaces

import (
	"encoding/json"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/services/fees/application"
	"github.com/leetatech/leeta_backend/services/fees/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/web"
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
// @Summary List fees.
// @Description The endpoint to get all list fees. Use filter t filter by type
// @Tags Fees
// @Accept json
// @produce json
// @param query.ResultSelector body query.ResultSelector true "list fees request body"
// @Security BearerToken
// @success 200 {object} query.ResponseListWithMetadata[models.Fee]
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /fees/ [PUT]
func (handler *FeesHttpHandler) FetchFeesHandler(w http.ResponseWriter, r *http.Request) {
	resultSelector, err := web.PrepareResultSelector(r, listFeesOptions, allowedSortFields, web.ResultSelectorDefaults(defaultSortingRequest))
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusBadRequest, leetError.ErrorResponseBody(leetError.InvalidRequestError, err))
		return
	}

	fees, totalRecord, err := handler.FeesApplication.GetTypedFees(r.Context(), resultSelector)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, leetError.ErrorResponseBody(leetError.InternalError, err))
		return
	}

	response := query.ResponseListWithMetadata[models.Fee]{
		Metadata: query.NewMetadata(resultSelector, totalRecord),
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
