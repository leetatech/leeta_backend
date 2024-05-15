package interfaces

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/services/auth/application"
	"github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/rs/zerolog/log"
	"net/http"
)

type AuthHttpHandler struct {
	AuthApplication application.AuthApplication
}

func NewAuthHttpHandler(authApplication application.AuthApplication) *AuthHttpHandler {
	return &AuthHttpHandler{
		AuthApplication: authApplication,
	}
}

// SignUpHandler godoc
// @Summary User Sign Up
// @Description The endpoint allows users, both vendors and customers to sign up
// @Tags Authentication
// @Accept json
// @Produce json
// @Param domain.SignupRequest body domain.SignupRequest true "user sign up request body"
// @Success 200 {object} domain.DefaultSigningResponse
// @Router /session/signup [post]
func (handler *AuthHttpHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var signUpRequest domain.SignupRequest
	err := json.NewDecoder(r.Body).Decode(&signUpRequest)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	token, err := handler.AuthApplication.SignUp(r.Context(), signUpRequest)
	if err != nil {
		log.Debug().Err(err).Msg("error completing user registration")
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	pkg.EncodeResult(w, token, http.StatusOK)
}

// RequestOTPHandler godoc
// @Summary Request for new OTP for target email
// @Description The endpoint allows client side to request for new OTP for target
// @Tags OTP Management
// @Accept json
// @Produce json
// @Param domain.EmailRequestBody body domain.EmailRequestBody true "request otp body"
// @Success 200 {object} pkg.DefaultResponse
// @Router /session/otp/request [post]
func (handler *AuthHttpHandler) RequestOTPHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.EmailRequestBody
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	token, err := handler.AuthApplication.RequestOTP(r.Context(), request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	token.Message = "OTP sent successfully"
	pkg.EncodeResult(w, token, http.StatusOK)
}

// EarlyAccessHandler godoc
// @Summary Early Access
// @Description The endpoint allows users to request for early access
// @Tags Early Access
// @Accept json
// @Produce json
// @Param models.EarlyAccess body models.EarlyAccess true "request early access body"
// @Success 200 {object} pkg.DefaultResponse
// @Router /session/early_access [post]
func (handler *AuthHttpHandler) EarlyAccessHandler(w http.ResponseWriter, r *http.Request) {
	var request models.EarlyAccess
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response, err := handler.AuthApplication.EarlyAccess(r.Context(), request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	pkg.EncodeResult(w, response, http.StatusOK)
}

// SignInHandler godoc
// @Summary User Sign In
// @Description The endpoint allows users, both vendors and customers to sign in
// @Tags Authentication
// @Accept json
// @Produce json
// @Param domain.SigningRequest body domain.SigningRequest true "user sign in request body"
// @Success 200 {object} domain.DefaultSigningResponse
// @Router /session/signin [post]
func (handler *AuthHttpHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	var signInRequest domain.SigningRequest
	err := json.NewDecoder(r.Body).Decode(&signInRequest)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	token, err := handler.AuthApplication.SignIn(r.Context(), signInRequest)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusInternalServerError)
		return
	}
	pkg.EncodeResult(w, token, http.StatusOK)
}

// ForgotPasswordHandler godoc
// @Summary Forgot Password
// @Description The endpoint allows users to request for password reset
// @Tags Password Management
// @Accept json
// @Produce json
// @Param domain.EmailRequestBody body domain.EmailRequestBody true "request forgot password body"
// @Success 200 {object} pkg.DefaultResponse
// @Router /session/password/forgot [post]
func (handler *AuthHttpHandler) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.EmailRequestBody
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response, err := handler.AuthApplication.ForgotPassword(r.Context(), request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusInternalServerError)
		return
	}

	pkg.EncodeResult(w, response, http.StatusOK)
}

// ValidateOTPHandler godoc
// @Summary Validate OTP
// @Description The endpoint allows users to validate OTP
// @Tags OTP Management
// @Accept json
// @Produce json
// @Param domain.OTPValidationRequest body domain.OTPValidationRequest true "request otp validation body"
// @Success 200 {object} pkg.DefaultResponse
// @Router /session/otp/validate [post]
func (handler *AuthHttpHandler) ValidateOTPHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.OTPValidationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response, err := handler.AuthApplication.ValidateOTP(r.Context(), request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	pkg.EncodeResult(w, response, http.StatusOK)
}

// CreateNewPasswordHandler godoc
// @Summary Create Password
// @Description The endpoint allows users to create a new password.
// @Tags Password Management
// @Accept json
// @Produce json
// @Param domain.CreateNewPasswordRequest body domain.CreateNewPasswordRequest true "request reset password body"
// @Success 200 {object} domain.APIResponseWithoutToken
// @Router /session/password/create [post]
func (handler *AuthHttpHandler) CreateNewPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.CreateNewPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response, err := handler.AuthApplication.CreateNewPassword(r.Context(), request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	pkg.EncodeResult(w, response, http.StatusOK)
}

// AdminSignUpHandler godoc
// @Summary Admin Sign Up
// @Description The endpoint allows admins to sign up
// @Tags Admin
// @Accept json
// @Produce json
// @Param domain.AdminSignUpRequest body domain.AdminSignUpRequest true "admin sign up request body"
// @Success 200 {object} domain.DefaultSigningResponse
// @Router /session/admin/signup [post]
func (handler *AuthHttpHandler) AdminSignUpHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.AdminSignUpRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	token, err := handler.AuthApplication.AdminSignUp(r.Context(), request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	pkg.EncodeResult(w, token, http.StatusOK)

}

// ReceiveGuestTokenHandler godoc
// @Summary Request accept guests
// @Description The endpoint to allow guests to shop
// @Tags Guest Management
// @Accept json
// @Produce json
// @Param domain.ReceiveGuestRequest body domain.ReceiveGuestRequest true "receive guest request body"
// @Success 200 {object} domain.ReceiveGuestResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /session/guest [post]
func (handler *AuthHttpHandler) ReceiveGuestTokenHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.ReceiveGuestRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusBadRequest, err)
		return
	}

	token, err := handler.AuthApplication.ReceiveGuestToken(request)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
		return
	}
	pkg.EncodeResult(w, token, http.StatusOK)
}

// UpdateGuestRecordHandler godoc
// @Summary Update guest record
// @Description The endpoint to update guest record
// @Tags Guest Management
// @Accept json
// @Produce json
// @Param models.Guest body models.Guest true "update guest request body"
// @Success 200 {object} pkg.DefaultResponse
// @error 400 {object} pkg.DefaultErrorResponse
// @error 401 {object} pkg.DefaultErrorResponse
// @Router /session/guest [put]
func (handler *AuthHttpHandler) UpdateGuestRecordHandler(w http.ResponseWriter, r *http.Request) {
	var request models.Guest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusBadRequest, err)
		return
	}

	resp, err := handler.AuthApplication.UpdateGuestRecord(r.Context(), request)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
		return
	}
	pkg.EncodeResult(w, resp, http.StatusOK)
}

// GetGuestRecordHandler godoc
// @Summary Get guest record
// @Description The endpoint to get guest record
// @Tags Guest Management
// @Accept json
// @Produce json
// @Param device_id path string true "device id"
// @Success 200 {object} models.Guest
// @error 400 {object} pkg.DefaultErrorResponse
// @error 401 {object} pkg.DefaultErrorResponse
// @Router /session/guest/{device_id} [get]
func (handler *AuthHttpHandler) GetGuestRecordHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "device_id")

	resp, err := handler.AuthApplication.GetGuestRecord(r.Context(), deviceID)
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusInternalServerError, err)
		return
	}

	pkg.EncodeResult(w, resp, http.StatusOK)
}
