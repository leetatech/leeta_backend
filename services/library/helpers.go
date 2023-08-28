package library

import (
	"encoding/base64"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"net/http"
)

func EncodeString(s string) string {
	data := base64.StdEncoding.EncodeToString([]byte(s))
	return string(data)
}

func DecodeString(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func CheckErrorType(err error, w http.ResponseWriter) {
	switch err := err.(type) {
	case *leetError.ErrorResponse:
		if err.Code() == leetError.ErrorUnauthorized {
			EncodeResult(w, err, http.StatusUnauthorized)
			return
		}
	default:
		EncodeResult(w, err, http.StatusInternalServerError)
		return
	}

}
