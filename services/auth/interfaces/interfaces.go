package interfaces

import (
	"encoding/json"
	"github.com/leetatech/leeta_backend/services/auth/application"
	"github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/library/models"
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
// @Param domain.SignUpRequest body domain.SignUpRequest true "user sign up request body"
// @Success 200 {object} domain.DefaultSigningResponse
// @Router /session/signup [post]
func (handler *AuthHttpHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var signUpRequest domain.SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&signUpRequest)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	token, err := handler.AuthApplication.SignUp(signUpRequest)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	library.EncodeResult(w, token, http.StatusOK)
}

// CreateOTPHandler godoc
// @Summary OTP Generation
// @Description The endpoint allows the generation of OTP
// @Tags OTP
// @Accept json
// @Produce json
// @Param domain.OTPRequest body domain.OTPRequest true "request otp body"
// @Success 200 {object} library.DefaultResponse
// @Router /session/otp/request [post]
func (handler *AuthHttpHandler) CreateOTPHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.OTPRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	token, err := handler.AuthApplication.CreateOTP(request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	library.EncodeResult(w, token, http.StatusOK)
}

// EarlyAccessHandler godoc
// @Summary Early Access
// @Description The endpoint allows users to request for early access
// @Tags EarlyAccess
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

	response, err := handler.AuthApplication.EarlyAccess(request)
	if err != nil {
		library.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	library.EncodeResult(w, response, http.StatusOK)
}
