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
// @Accept multipart/form-data
// @Produce json
// @Param first_name formData string true "First name of the vendor"
// @Param last_name formData string true "Last name of the vendor"
// @Param business_name formData string true "Business name of the vendor"
// @Param cac formData string true "CAC number of the vendor"
// @Param business_category formData string true "Business category of the vendor"
// @Param description formData string true "Description of the vendor"
// @Param primary_phone formData bool true "Is the primary phone number"
// @Param phone_number formData string true "Phone number of the vendor"
// @Param state formData string true "State of the vendor"
// @Param city formData string true "City of the vendor"
// @Param lga formData string true "Local Government Area of the vendor"
// @Param full_address formData string true "Full address of the vendor"
// @Param closest_landmark formData string true "Closest landmark to the vendor's location"
// @Param latitude formData string true "Latitude of the vendor's location"
// @Param longitude formData string true "Longitude of the vendor's location"
// @Param image formData file true "Image of the vendor"
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
		library.CheckErrorType(err, w)
		return
	}
	library.EncodeResult(w, token, http.StatusOK)
}

// AddVendorByAdminHandler godoc
// @Summary Admin adds vendor and business
// @Description The endpoint allows the admin to add a vendor and their business
// @Tags user/admin/vendor
// @Accept multipart/form-data
// @Produce json
// @Param first_name formData string true "First name of the vendor"
// @Param last_name formData string true "Last name of the vendor"
// @Param business_name formData string true "Business name of the vendor"
// @Param cac formData string true "CAC number of the vendor"
// @Param business_category formData string true "Business category of the vendor"
// @Param description formData string true "Description of the vendor"
// @Param primary_phone formData bool true "Is the primary phone number"
// @Param phone_number formData string true "Phone number of the vendor"
// @Param state formData string true "State of the vendor"
// @Param city formData string true "City of the vendor"
// @Param lga formData string true "Local Government Area of the vendor"
// @Param full_address formData string true "Full address of the vendor"
// @Param closest_landmark formData string true "Closest landmark to the vendor's location"
// @Param latitude formData string true "Latitude of the vendor's location"
// @Param longitude formData string true "Longitude of the vendor's location"
// @Param image formData file true "Image of the vendor"
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
		library.CheckErrorType(err, w)
		return
	}
	library.EncodeResult(w, token, http.StatusOK)
}
