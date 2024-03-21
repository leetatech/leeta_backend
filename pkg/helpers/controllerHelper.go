package helpers

import (
	"errors"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/filter"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"net/http"
)

func CheckErrorType(err error, w http.ResponseWriter) {
	var lerr *leetError.ErrorResponse
	switch {
	case errors.As(err, &lerr):
		if lerr.Code() == leetError.ErrorUnauthorized {
			pkg.EncodeErrorResult(w, http.StatusUnauthorized)
			return
		}
	default:
		pkg.EncodeErrorResult(w, http.StatusInternalServerError)
		return
	}

}

func ValidateQueryFilter(request *filter.PagingRequest) (*filter.PagingRequest, error) {
	if request == nil {
		return nil, leetError.ErrorResponseBody(leetError.InvalidPageRequestError, errors.New("the paging field is required but it is missing"))
	}

	if request.PageIndex == 0 {
		request.PageIndex = 1
	}

	if request.PageSize == 0 {
		request.PageSize = 10
	}

	return request, nil
}
