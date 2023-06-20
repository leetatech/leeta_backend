package domain

import "github.com/leetatech/leeta_backend/services/library/models"

type SignUpRequest struct {
	Email    string              `json:"email"`
	Password string              `json:"password"`
	UserType models.UserCategory `json:"user_type"`
} // @name SignUpRequest

type DefaultSigningResponse struct {
	AuthToken string `json:"auth_token,omitempty"`
} // @name DefaultSigningResponse
