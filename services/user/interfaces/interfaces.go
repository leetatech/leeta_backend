package interfaces

import (
	"encoding/json"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/user/application"
	"github.com/leetatech/leeta_backend/services/user/domain"
	"net/http"
)

type UserHttpHandler struct {
	UserApplication application.UserApplication
}

func NewUserHttpHandler(userApplication application.UserApplication) *UserHttpHandler {
	return &UserHttpHandler{
		UserApplication: userApplication,
	}
}

// VendorVerificationHandler godoc
// @Summary Vendor Verification
// @Description The endpoint allows the verification process of vendor
// @Tags user/vendor
// @Accept json
// @Produce json
// @Param domain.VendorVerificationRequest body domain.VendorVerificationRequest true "vendor verification request body"
// @Security BearerToken
// @Success 200 {object} library.DefaultResponse
// @Failure 401 {object} library.DefaultErrorResponse
// @Failure 400 {object} library.DefaultErrorResponse
// @Router /user/vendor/verification [post]
func (handler *UserHttpHandler) VendorVerificationHandler(w http.ResponseWriter, r *http.Request) {
	var verificationRequest domain.VendorVerificationRequest
	err := json.NewDecoder(r.Body).Decode(&verificationRequest)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	token, err := handler.UserApplication.VendorVerification(r.Context(), verificationRequest)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	library.EncodeResult(w, token, http.StatusOK)
}
