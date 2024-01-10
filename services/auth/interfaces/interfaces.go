package interfaces

import (
	"encoding/json"
	"github.com/leetatech/leeta_backend/services/auth/application"
	"github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/library/models"
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
// @Description The endpoint allows users, both vendors and buyers to sign up
// @Tags Session
// @Accept json
// @Produce json
// @Param domain.SignupRequest body domain.SignupRequest true "user sign up request body"
// @Success 200 {object} domain.DefaultSigningResponse
// @Router /session/signup [post]
func (handler *AuthHttpHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var signUpRequest domain.SignupRequest
	err := json.NewDecoder(r.Body).Decode(&signUpRequest)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	token, err := handler.AuthApplication.SignUp(r.Context(), signUpRequest)
	if err != nil {
		log.Debug().Err(err).Msg("error completing user registration")
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	library.EncodeResult(w, token, http.StatusOK)
}

// RequestOTPHandler godoc
// @Summary Request for new OTP for target email
// @Description The endpoint allows client side to request for new OTP for target
// @Tags Session
// @Accept json
// @Produce json
// @Param domain.EmailRequestBody body domain.EmailRequestBody true "request otp body"
// @Success 200 {object} library.DefaultResponse
// @Router /session/otp/request [post]
func (handler *AuthHttpHandler) RequestOTPHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.EmailRequestBody
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	token, err := handler.AuthApplication.RequestOTP(r.Context(), request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	token.Message = "OTP sent successfully"
	library.EncodeResult(w, token, http.StatusOK)
}

// EarlyAccessHandler godoc
// @Summary Early Access
// @Description The endpoint allows users to request for early access
// @Tags Early Access
// @Accept json
// @Produce json
// @Param models.EarlyAccess body models.EarlyAccess true "request early access body"
// @Success 200 {object} library.DefaultResponse
// @Router /session/early_access [post]
func (handler *AuthHttpHandler) EarlyAccessHandler(w http.ResponseWriter, r *http.Request) {
	var request models.EarlyAccess
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response, err := handler.AuthApplication.EarlyAccess(r.Context(), request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	library.EncodeResult(w, response, http.StatusOK)
}

// SignInHandler godoc
// @Summary User Sign In
// @Description The endpoint allows users, both vendors and buyers to sign in
// @Tags Session
// @Accept json
// @Produce json
// @Param domain.SigningRequest body domain.SigningRequest true "user sign in request body"
// @Success 200 {object} domain.DefaultSigningResponse
// @Router /session/signin [post]
func (handler *AuthHttpHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	var signInRequest domain.SigningRequest
	err := json.NewDecoder(r.Body).Decode(&signInRequest)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	token, err := handler.AuthApplication.SignIn(r.Context(), signInRequest)
	if err != nil {
		library.EncodeResult(w, err, http.StatusInternalServerError)
		return
	}
	library.EncodeResult(w, token, http.StatusOK)
}

// ForgotPasswordHandler godoc
// @Summary Forgot Password
// @Description The endpoint allows users to request for password reset
// @Tags Session
// @Accept json
// @Produce json
// @Param domain.EmailRequestBody body domain.EmailRequestBody true "request forgot password body"
// @Success 200 {object} library.DefaultResponse
// @Router /session/forgot_password [post]
func (handler *AuthHttpHandler) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.EmailRequestBody
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response, err := handler.AuthApplication.ForgotPassword(r.Context(), request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusInternalServerError)
		return
	}

	library.EncodeResult(w, response, http.StatusOK)
}

// ValidateOTPHandler godoc
// @Summary Validate OTP
// @Description The endpoint allows users to validate OTP
// @Tags Session
// @Accept json
// @Produce json
// @Param domain.OTPValidationRequest body domain.OTPValidationRequest true "request otp validation body"
// @Success 200 {object} library.DefaultResponse
// @Router /session/otp/validate [post]
func (handler *AuthHttpHandler) ValidateOTPHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.OTPValidationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response, err := handler.AuthApplication.ValidateOTP(r.Context(), request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	library.EncodeResult(w, response, http.StatusOK)
}

// ResetPasswordHandler godoc
// @Summary Reset Password
// @Description The endpoint allows users to reset password
// @Tags Session
// @Accept json
// @Produce json
// @Param domain.ResetPasswordRequest body domain.ResetPasswordRequest true "request reset password body"
// @Success 200 {object} domain.DefaultSigningResponse
// @Router /session/reset_password [post]
func (handler *AuthHttpHandler) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.ResetPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	response, err := handler.AuthApplication.ResetPassword(r.Context(), request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	library.EncodeResult(w, response, http.StatusOK)
}

// AdminSignUpHandler godoc
// @Summary Admin Sign Up
// @Description The endpoint allows admins to sign up
// @Tags Session
// @Accept json
// @Produce json
// @Param domain.AdminSignUpRequest body domain.AdminSignUpRequest true "admin sign up request body"
// @Success 200 {object} domain.DefaultSigningResponse
// @Router /session/admin/signup [post]
func (handler *AuthHttpHandler) AdminSignUpHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.AdminSignUpRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	token, err := handler.AuthApplication.AdminSignUp(r.Context(), request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	library.EncodeResult(w, token, http.StatusOK)

}
