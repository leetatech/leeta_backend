package leetError

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

type ErrorCode int

const (
	DatabaseError              ErrorCode = 1001
	DatabaseNoRecordError      ErrorCode = 1002
	UnmarshalError             ErrorCode = 1003
	MarshalError               ErrorCode = 1004
	PasswordValidationError    ErrorCode = 1005
	EncryptionError            ErrorCode = 1006
	DecryptionError            ErrorCode = 1007
	DuplicateUserError         ErrorCode = 1008
	UserNotFoundError          ErrorCode = 1009
	IdentityNotFoundError      ErrorCode = 1010
	UserLockedError            ErrorCode = 1011
	CredentialsValidationError ErrorCode = 1012
	TokenGenerationError       ErrorCode = 1013
	TokenValidationError       ErrorCode = 1014
	UserCategoryError          ErrorCode = 1015
	EmailSendingError          ErrorCode = 1016
	BusinessCategoryError      ErrorCode = 1017
	StatusesError              ErrorCode = 1018
	ErrorUnauthorized          ErrorCode = 1019
	EmailFormatError           ErrorCode = 1020
	ValidEmailHostError        ErrorCode = 1021
	ValidLeetaDomainError      ErrorCode = 1022
	FormParseError             ErrorCode = 1023
)

var (
	errorTypes = map[ErrorCode]string{
		DatabaseError:              "DatabaseError",
		DatabaseNoRecordError:      "DatabaseNoRecordError",
		UnmarshalError:             "UnmarshalError",
		MarshalError:               "MarshalError",
		PasswordValidationError:    "PasswordValidationError",
		EncryptionError:            "EncryptionError",
		DecryptionError:            "DecryptionError",
		DuplicateUserError:         "DuplicateUserError",
		UserNotFoundError:          "UserNotFoundError",
		IdentityNotFoundError:      "IdentityNotFoundError",
		UserLockedError:            "UserLockedError",
		CredentialsValidationError: "CredentialsValidationError",
		TokenGenerationError:       "TokenGenerationError",
		TokenValidationError:       "TokenValidationError",
		UserCategoryError:          "UserCategoryError",
		EmailSendingError:          "EmailSendingError",
		BusinessCategoryError:      "BusinessCategoryError",
		StatusesError:              "StatusesError",
		ErrorUnauthorized:          "ErrorUnauthorized",
		EmailFormatError:           "EmailFormatError",
		ValidEmailHostError:        "ValidEmailHostError",
		ValidLeetaDomainError:      "ValidLeetaDomainError",
		FormParseError:             "FormParseError",
	}

	errorMessages = map[ErrorCode]string{
		DatabaseError:              "An error occurred while reading from the database",
		DatabaseNoRecordError:      "An error occurred because no record was found",
		UnmarshalError:             "An error occurred while unmarshalling data",
		MarshalError:               "An error occurred while marshaling data",
		PasswordValidationError:    "An error occurred while validating password. | Password must contain at least six character long, one uppercase letter, one lowercase letter, one digit, and one special character | password and confirm password don't match",
		EncryptionError:            "An error occurred while encrypting",
		DecryptionError:            "An error occurred while decrypting",
		DuplicateUserError:         "An error occurred because user already exists",
		UserNotFoundError:          "An error occurred because this is not a registered user",
		IdentityNotFoundError:      "An error occurred because this is not a registered identity",
		UserLockedError:            "An error occurred because this user is locked",
		CredentialsValidationError: "An error occurred because the credentials are invalid",
		TokenGenerationError:       "An error occurred while generating token",
		TokenValidationError:       "An error occurred because the token is invalid | validated | expired",
		UserCategoryError:          "An error occurred because the user category is invalid",
		EmailSendingError:          "An error occurred while sending email",
		BusinessCategoryError:      "An error occurred because the business category is invalid",
		StatusesError:              "An error occurred because the statuses are invalid",
		ErrorUnauthorized:          "An error occurred because the user is unauthorized",
		EmailFormatError:           "An error occurred because the email format is invalid",
		ValidEmailHostError:        "An error occurred because the domain does not exist or cannot receive emails",
		ValidLeetaDomainError:      "An error occurred because the domain does not belong to leeta or cannot receive emails",
		FormParseError:             "An error occurred because the form parse failed or file retrieval failed",
	}
)

type ErrorResponse struct {
	ErrorReference uuid.UUID `json:"error_reference"`
	Code           ErrorCode `json:"code"`
	ErrorType      string    `json:"error_type"`
	Message        string    `json:"message"`
	Err            error     `json:"-"`
	StackTrace     string    `json:"-"`
	TimeStamp      string    `json:"timestamp"`
}

func (err ErrorResponse) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return err.Message
}

func (err ErrorResponse) Unwrap() error {
	return err.Err
}

func (err ErrorResponse) Wrap(message string) error {
	return errors.Wrap(err.Err, message)
}

func ErrorResponseBody(code ErrorCode, err error) error {
	errorResponse := ErrorResponse{
		ErrorReference: uuid.New(),
		Code:           code,
		ErrorType:      errorTypes[code],
		Message:        errorMessages[code],
		Err:            err,
		TimeStamp:      time.Now().Format(time.RFC3339),
	}

	// Capture stack trace if available
	if err != nil {
		errorResponse.StackTrace = fmt.Sprintf("%+v", errors.WithStack(err).Error())
	}

	return errorResponse
}

func ErrorMessage(code ErrorCode) string {
	return errorMessages[code]
}

func ErrorType(code ErrorCode) string {
	return errorTypes[code]
}
