package interfaces

import (
	"encoding/json"
	"net/http"

	_ "github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/helpers"
	"github.com/leetatech/leeta_backend/pkg/jwtmiddleware"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/user/application"
)

type UserHttpHandler struct {
	UserApplication application.UserApplication
}

func New(userApplication application.UserApplication) *UserHttpHandler {
	return &UserHttpHandler{
		UserApplication: userApplication,
	}
}

// VendorVerificationHandler godoc
// @Summary Vendor Verification
// @Description The endpoint allows the verification process of vendor
// @Tags Vendor
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
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /user/vendor/verification [post]
func (handler *UserHttpHandler) VendorVerificationHandler(w http.ResponseWriter, r *http.Request) {
	request, err := checkFormFileSpecification(r)
	if err != nil {
		return
	}

	token, err := handler.UserApplication.VendorVerification(r.Context(), *request)
	if err != nil {
		helpers.CheckErrorType(err, w)
		return
	}
	jwtmiddleware.WriteJSONResponse(w, token, http.StatusOK)
}

// AddVendorByAdminHandler godoc
// @Summary Admin adds vendor and business
// @Description The endpoint allows the admin to add a vendor and their business
// @Tags Admin
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
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /user/admin/vendor [post]
func (handler *UserHttpHandler) AddVendorByAdminHandler(w http.ResponseWriter, r *http.Request) {
	request, err := checkFormFileSpecification(r)
	if err != nil {
		return
	}

	token, err := handler.UserApplication.AddVendorByAdmin(r.Context(), *request)
	if err != nil {
		helpers.CheckErrorType(err, w)
		return
	}
	jwtmiddleware.WriteJSONResponse(w, token, http.StatusOK)
}

// UpdateUserData godoc
// @Summary Update User data
// @Description Update user data is the endpoint used to make changes to a user database record
// @Tags User
// @Accept json
// @Produce json
// @Param models.User body models.User true "update user record"
// @Security BearerToken
// @Success 204 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /user/ [put]
func (handler *UserHttpHandler) UpdateUserData(w http.ResponseWriter, r *http.Request) {
	var request models.User

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		jwtmiddleware.WriteJSONErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	resp, err := handler.UserApplication.UpdateRecord(r.Context(), request)
	if err != nil {
		helpers.CheckErrorType(err, w)
		return
	}

	jwtmiddleware.WriteJSONResponse(w, resp, http.StatusOK)
}

// Data godoc
// @Summary Get authenticated user data
// @Description The endpoint to get user record from current user jwt token
// @Tags User
// @Produce json
// @Security BearerToken
// @Success 200 {object} models.Customer
// @error 400 {object} pkg.DefaultErrorResponse
// @error 401 {object} pkg.DefaultErrorResponse
// @Router /user/ [get]
func (handler *UserHttpHandler) Data(w http.ResponseWriter, r *http.Request) {
	resp, err := handler.UserApplication.Data(r.Context())
	if err != nil {
		helpers.CheckErrorType(err, w)
		return
	}

	jwtmiddleware.WriteJSONResponse(w, resp, http.StatusOK)
}
