package helpers

import (
	"errors"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/query"
	"net/http"
)

func CheckErrorType(err error, w http.ResponseWriter) {
	var lerr *leetError.ErrorResponse
	switch {
	case errors.As(err, &lerr):
		if lerr.Code == leetError.ErrorUnauthorized {
			pkg.EncodeErrorResult(w, http.StatusUnauthorized, err)
			return
		}
	default:
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
		return
	}

}

func ValidateResultSelector(resultSelector *query.ResultSelector) (*query.ResultSelector, error) {
	if resultSelector == nil {
		return nil, leetError.ErrorResponseBody(leetError.InvalidRequestError, errors.New("the result selector cannot be empty"))
	}

	if resultSelector.Paging == nil {
		return nil, leetError.ErrorResponseBody(leetError.InvalidRequestError, errors.New("the paging cannot be empty"))
	}

	if err := resultSelector.Paging.Validate(); err != nil {
		return nil, leetError.ErrorResponseBody(leetError.InvalidRequestError, fmt.Errorf("invalid paging request %w", err))
	}
	resultSelector.Paging.ApplyDefaults()

	return resultSelector, nil
}
