package interfaces

import (
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/user/application"
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
	request, err := checkFormFileSpecification(r)
	if err != nil {
		return
	}

	token, err := handler.UserApplication.VendorVerification(r.Context(), *request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	library.EncodeResult(w, token, http.StatusOK)
}

// AddVendorByAdminHandler godoc
// @Summary Admin adds vendor and business
// @Description The endpoint allows the admin to add a vendor and their business
// @Tags user/admin/vendor
// @Accept json
// @Produce json
// @Param domain.VendorVerificationRequest body domain.VendorVerificationRequest true "vendor verification request body"
// @Security BearerToken
// @Success 200 {object} library.DefaultResponse
// @Failure 401 {object} library.DefaultErrorResponse
// @Failure 400 {object} library.DefaultErrorResponse
// @Router /user/admin/vendor [post]
func (handler *UserHttpHandler) AddVendorByAdminHandler(w http.ResponseWriter, r *http.Request) {
	request, err := checkFormFileSpecification(r)
	if err != nil {
		return
	}

	token, err := handler.UserApplication.AddVendorByAdmin(r.Context(), *request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	library.EncodeResult(w, token, http.StatusOK)
}
