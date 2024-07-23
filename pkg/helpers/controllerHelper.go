package helpers

import (
	"errors"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
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
		case leetError.DatabaseNoRecordError, leetError.LGANotFoundError:
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
