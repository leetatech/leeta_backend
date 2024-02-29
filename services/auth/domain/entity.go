package domain

import "github.com/leetatech/leeta_backend/services/library/models"

type SignupRequest struct {
	FullName string              `json:"full_name"`
	Email    string              `json:"email"`
	Password string              `json:"password"`
	UserType models.UserCategory `json:"user_type"`
} // @name SignupRequest

type SigningRequest struct {
	Email    string              `json:"email"`
	Password string              `json:"password"`
	UserType models.UserCategory `json:"user_type"`
} // @name SigningRequest

type DefaultSigningResponse struct {
	AuthToken string `json:"auth_token,omitempty"`
	Body      any    `json:"body"`
} // @name DefaultSigningResponse

type APIResponseWithoutToken struct {
	Body any `json:"body"`
} // @name APIResponseWithoutToken

type OTPRequest struct {
	Topic        string                     `json:"topic" bson:"topic"`
	Type         models.MessageDeliveryType `json:"type" bson:"type"`
	Target       string                     `json:"target" bson:"target"`
	UserCategory models.UserCategory        `json:"userCategory" bson:"user_category"`
} // @name OTPRequest

type EmailRequestBody struct {
	Email string `json:"email" bson:"email"`
} // @name EmailRequestBody

type OTPValidationRequest struct {
	Code   string `json:"code" bson:"code"`
	Target string `json:"target" bson:"target"`
} // @name OTPValidationRequest

type CreateNewPasswordRequest struct {
	Email           string `json:"email" bson:"email"`
	Password        string `json:"password" bson:"password"`
	ConfirmPassword string `json:"confirm_password" bson:"confirm_password"`
} // @name CreateNewPasswordRequest

type AdminSignUpRequest struct {
	Email      string         `json:"email"`
	Password   string         `json:"password"`
	FirstName  string         `json:"first_name"`
	LastName   string         `json:"last_name"`
	DOB        string         `json:"dob" bson:"dob"`
	Address    models.Address `json:"address" bson:"address"`
	Phone      string         `json:"phone" bson:"phone"`
	Department string         `json:"department"`
	Role       string         `json:"role"`
} // @name AdminSignUpRequest

type ReceiveGuestRequest struct {
	Guest bool `json:"guest" bson:"guest"`
} // @name ReceiveGuestRequest

type ReceiveGuestResponse struct {
	SessionID string `json:"session_id" bson:"session_id"`
	Token     string `json:"token" bson:"token"`
} // @name ReceiveGuestResponse
