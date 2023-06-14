package interfaces

import (
	"encoding/json"
	"github.com/leetatech/leeta_backend/services/order/application"
	"net/http"
)

type HTTPHandler struct {
	OrderApplication application.OrderApplication
}

func NewOrderHTTPHandler(orderApplication application.OrderApplication) *HTTPHandler {
	return &HTTPHandler{
		OrderApplication: orderApplication,
	}

}

func encodeResult(w http.ResponseWriter, result interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	data := struct {
		Data interface{} `json:"data"`
	}{
		Data: result,
	}

	err := json.NewEncoder(w).Encode(&data)
	if err != nil {
		return
	}
}
