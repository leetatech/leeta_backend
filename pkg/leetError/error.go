package leetError

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"runtime"
	"time"
)

type ErrorCode int

const (
	DatabaseError                ErrorCode = 1001
	DatabaseNoRecordError        ErrorCode = 1002
	UnmarshalError               ErrorCode = 1003
	MarshalError                 ErrorCode = 1004
	PasswordValidationError      ErrorCode = 1005
	EncryptionError              ErrorCode = 1006
	DecryptionError              ErrorCode = 1007
	DuplicateUserError           ErrorCode = 1008
	UserNotFoundError            ErrorCode = 1009
	IdentityNotFoundError        ErrorCode = 1010
	UserLockedError              ErrorCode = 1011
	CredentialsValidationError   ErrorCode = 1012
	TokenGenerationError         ErrorCode = 1013
	TokenValidationError         ErrorCode = 1014
	UserCategoryError            ErrorCode = 1015
	EmailSendingError            ErrorCode = 1016
	BusinessCategoryError        ErrorCode = 1017
	StatusesError                ErrorCode = 1018
	ErrorUnauthorized            ErrorCode = 1019
	EmailFormatError             ErrorCode = 1020
	ValidEmailHostError          ErrorCode = 1021
	ValidLeetaDomainError        ErrorCode = 1022
	FormParseError               ErrorCode = 1023
	OrderStatusesError           ErrorCode = 1024
	ProductCategoryError         ErrorCode = 1025
	ProductSubCategoryError      ErrorCode = 1026
	ProductStatusError           ErrorCode = 1027
	ForgotPasswordError          ErrorCode = 1028
	MissingUserNames             ErrorCode = 1029
	InvalidUserRoleError         ErrorCode = 1030
	InvalidIdentityError         ErrorCode = 1031
	InvalidOTPError              ErrorCode = 1032
	CartStatusesError            ErrorCode = 1033
	AmountPaidError              ErrorCode = 1034
	FeesStatusesError            ErrorCode = 1035
	InvalidPageRequestError      ErrorCode = 1036
	CartItemQuantityError        ErrorCode = 1037
	CartItemRequestQuantityError ErrorCode = 1038
	InvalidRequestError          ErrorCode = 1039 // generic
	InternalError                ErrorCode = 1040
	InvalidProductIdError        ErrorCode = 1041
)

var (
	errorTypes = map[ErrorCode]string{
		DatabaseError:                "DatabaseError",
		DatabaseNoRecordError:        "DatabaseNoRecordError",
		UnmarshalError:               "UnmarshalError",
		MarshalError:                 "MarshalError",
		PasswordValidationError:      "PasswordValidationError",
		EncryptionError:              "EncryptionError",
		DecryptionError:              "DecryptionError",
		DuplicateUserError:           "DuplicateUserError",
		UserNotFoundError:            "UserNotFoundError",
		IdentityNotFoundError:        "IdentityNotFoundError",
		UserLockedError:              "UserLockedError",
		CredentialsValidationError:   "CredentialsValidationError",
		TokenGenerationError:         "TokenGenerationError",
		TokenValidationError:         "TokenValidationError",
		UserCategoryError:            "UserCategoryError",
		EmailSendingError:            "EmailSendingError",
		BusinessCategoryError:        "BusinessCategoryError",
		StatusesError:                "StatusesError",
		ErrorUnauthorized:            "ErrorUnauthorized",
		EmailFormatError:             "EmailFormatError",
		ValidEmailHostError:          "ValidEmailHostError",
		ValidLeetaDomainError:        "ValidLeetaDomainError",
		FormParseError:               "FormParseError",
		OrderStatusesError:           "OrderStatusesError",
		ProductCategoryError:         "ProductCategoryError",
		ProductSubCategoryError:      "ProductSubCategoryError",
		ProductStatusError:           "ProductStatusError",
		ForgotPasswordError:          "ForgotPasswordError",
		MissingUserNames:             "MissingUserNamesError",
		InvalidUserRoleError:         "InvalidUserRoleError",
		InvalidIdentityError:         "InvalidIdentityError",
		InvalidOTPError:              "InvalidOTPError",
		CartStatusesError:            "CartStatusesError",
		AmountPaidError:              "AmountPaidError",
		FeesStatusesError:            "FeesStatusesError",
		InvalidPageRequestError:      "InvalidPageRequestError",
		CartItemQuantityError:        "CartItemQuantityError",
		CartItemRequestQuantityError: "CartItemRequestQuantityError",
		InvalidRequestError:          "InvalidRequestError",
		InternalError:                "InternalError",
		InvalidProductIdError:        "InvalidProductIdError",
	}

	errorMessages = map[ErrorCode]string{
		DatabaseError:                "An error occurred while reading from the database",
		DatabaseNoRecordError:        "An error occurred because no record was found",
		UnmarshalError:               "An error occurred while unmarshalling data",
		MarshalError:                 "An error occurred while marshaling data",
		PasswordValidationError:      "An error occurred while validating password. | Password must contain at least six character long, one uppercase letter, one lowercase letter, one digit, and one special character | password and confirm password don't match",
		EncryptionError:              "An error occurred while encrypting",
		DecryptionError:              "An error occurred while decrypting",
		DuplicateUserError:           "An error occurred because user already exists",
		UserNotFoundError:            "An error occurred because this is not a registered user",
		IdentityNotFoundError:        "An error occurred because this user identity is not known",
		UserLockedError:              "An error occurred because this user is locked",
		CredentialsValidationError:   "An error occurred because the credentials are invalid",
		TokenGenerationError:         "An error occurred while generating token",
		TokenValidationError:         "An error occurred because the token is invalid | validated | expired",
		UserCategoryError:            "An error occurred because the user category is invalid",
		EmailSendingError:            "An error occurred while sending email",
		BusinessCategoryError:        "An error occurred because the business category is invalid",
		StatusesError:                "An error occurred because the statuses are invalid",
		ErrorUnauthorized:            "An error occurred because the user is unauthorized",
		EmailFormatError:             "An error occurred because the email format is invalid",
		ValidEmailHostError:          "An error occurred because the domain does not exist or cannot receive emails",
		ValidLeetaDomainError:        "An error occurred because the domain does not belong to leeta or cannot receive emails",
		FormParseError:               "An error occurred because the form parse failed or file retrieval failed",
		OrderStatusesError:           "An error occurred because the order status is invalid",
		ProductCategoryError:         "An error occurred because the product category is invalid",
		ProductSubCategoryError:      "An error occurred because the product subcategory is invalid",
		ProductStatusError:           "An error occurred because the product status is invalid",
		ForgotPasswordError:          "An error occurred while trying to reset a user password",
		MissingUserNames:             "An error occurred because user first name/last name was not found",
		InvalidUserRoleError:         "An error occurred because the user is trying to login with the wrong app",
		InvalidIdentityError:         "An error occurred because the user identity data is invalid",
		InvalidOTPError:              "An error occurred because the OTP is invalid",
		CartStatusesError:            "An error occurred because the cart status is invalid",
		AmountPaidError:              "An error occurred because the amount paid is invalid",
		FeesStatusesError:            "An error occurred because the fees status is invalid",
		InvalidPageRequestError:      "An error occurred because the page request field is required",
		CartItemQuantityError:        "An error occurred because the stored cart item quantity/weight is already 0. Please delete the item or increase the quantity to continue",
		CartItemRequestQuantityError: "An error occurred because the request quantity/weight field is 0. Please increase the quantity/weight to continue",
		InvalidRequestError:          "An error occurred because the request is invalid",
		InternalError:                "An error has occurred in the server",
		InvalidProductIdError:        "An error occurred because the product id is invalid",
	}
)

type ErrorResponse struct {
	ErrorReference uuid.UUID `json:"error_reference"`
	ErrorCode      ErrorCode `json:"error_code"`
	Code           ErrorCode `json:"-"`
	ErrorType      string    `json:"error_type"`
	Message        string    `json:"message"`
	Err            any       `json:"internal_error_message"`
	StackTrace     string    `json:"-"`
	File           string    `json:"-"`
	Line           int       `json:"-"`
	TimeStamp      string    `json:"-"`
}

func (e ErrorResponse) Error() string {
	return e.Format()
}

func (e ErrorResponse) Format() string {
	return fmt.Sprintf("%s:%s | %s:%s | %s:%d | stackTrace:%s", e.ErrorReference, e.Err, e.ErrorType, e.Message, e.File, e.Line, e.StackTrace)
}

func ErrorResponseBody(code ErrorCode, err error) error {
	_, file, line, _ := runtime.Caller(1)
	errorResponse := ErrorResponse{
		ErrorReference: uuid.New(),
		ErrorCode:      code,
		ErrorType:      errorTypes[code],
		Message:        errorMessages[code],
		Err:            err.Error(),
		File:           file,
		Line:           line,
		TimeStamp:      time.Now().Format(time.RFC3339),
	}

	// Capture stack trace if available
	errorResponse.StackTrace = fmt.Sprintf("%+v", errors.WithStack(err).Error())

	return errorResponse
}

func ErrorMessage(code ErrorCode) string {
	return errorMessages[code]
}

func ErrorType(code ErrorCode) string {
	return errorTypes[code]
}
