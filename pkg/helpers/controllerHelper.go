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
	var lerr leetError.ErrorResponse
	switch {
	case errors.As(err, &lerr):
		switch lerr.ErrorCode {
		case leetError.ErrorUnauthorized:
			pkg.EncodeErrorResult(w, http.StatusUnauthorized, err)
			return
		case leetError.DatabaseNoRecordError:
			pkg.EncodeErrorResult(w, http.StatusNotFound, err)
			return
		case leetError.InvalidRequestError:
			pkg.EncodeErrorResult(w, http.StatusBadRequest, err)
			return
		default:
			pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
			return
		}
	default:
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
		return
	}

}

func ValidateResultSelector(resultSelector query.ResultSelector) (query.ResultSelector, error) {
	if resultSelector.Paging == nil {
		return resultSelector, leetError.ErrorResponseBody(leetError.InvalidRequestError, errors.New("the paging cannot be empty"))
	}

	if err := resultSelector.Paging.Validate(); err != nil {
		return resultSelector, leetError.ErrorResponseBody(leetError.InvalidRequestError, fmt.Errorf("invalid paging request %w", err))
	}
	resultSelector.Paging.ApplyDefaults()

	return resultSelector, nil
}
