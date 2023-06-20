package interfaces

import (
	"encoding/json"
	"github.com/leetatech/leeta_backend/services/auth/application"
	"github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/library"
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
		library.EncodeResult(w, err, http.StatusOK)
		return
	}

	token, err := handler.AuthApplication.SignUp(signUpRequest)
	if err != nil {
		library.EncodeResult(w, err, http.StatusOK)
		return
	}
	library.EncodeResult(w, token, http.StatusOK)
}
