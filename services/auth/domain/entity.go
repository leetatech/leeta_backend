package domain

import "github.com/leetatech/leeta_backend/services/library/models"

type SigningRequest struct {
	Email    string              `json:"email"`
	Password string              `json:"password"`
	UserType models.UserCategory `json:"user_type"`
} // @name SigningRequest

type DefaultSigningResponse struct {
	AuthToken string `json:"auth_token,omitempty"`
} // @name DefaultSigningResponse

type OTPRequest struct {
	Topic        string                     `json:"topic" bson:"topic"`
	Type         models.MessageDeliveryType `json:"type" bson:"type"`
	Target       string                     `json:"target" bson:"target"`
	UserCategory models.UserCategory        `json:"userCategory" bson:"user_category"`
} // @name OTPRequest
