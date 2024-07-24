package helpers

import (
	"errors"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/pkg/jwtmiddleware"
	"net/http"
)

func CheckErrorType(err error, w http.ResponseWriter) {
	var lerr *errs.Response
	switch {
	case errors.As(err, &lerr):
		switch lerr.ErrorCode {
		case errs.ErrorUnauthorized:
			jwtmiddleware.WriteJSONErrorResponse(w, http.StatusUnauthorized, err)
			return
		case errs.DatabaseNoRecordError, errs.LGANotFoundError:
			jwtmiddleware.WriteJSONErrorResponse(w, http.StatusNotFound, err)
			return
		case errs.InvalidRequestError:
			jwtmiddleware.WriteJSONErrorResponse(w, http.StatusBadRequest, err)
			return
		default:
			jwtmiddleware.WriteJSONErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	default:
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

}
