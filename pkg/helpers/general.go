package helpers

import (
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"net/http"
)

func CheckErrorType(err error, w http.ResponseWriter) {
	switch err := err.(type) {
	case *leetError.ErrorResponse:
		if err.Code() == leetError.ErrorUnauthorized {
			pkg.EncodeErrorResult(w, http.StatusUnauthorized)
			return
		}
	default:
		pkg.EncodeErrorResult(w, http.StatusInternalServerError)
		return
	}

}
