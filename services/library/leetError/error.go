package leetError

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

type ErrorCode int

const (
	DatabaseError         ErrorCode = 1001
	DatabaseNoRecordError ErrorCode = 1002
	UnmarshalError        ErrorCode = 1003
	MarshalError          ErrorCode = 1004
)

var (
	errorTypes = map[ErrorCode]string{
		DatabaseError:         "DatabaseError",
		DatabaseNoRecordError: "DatabaseNoRecordError",
		UnmarshalError:        "UnmarshalError",
		MarshalError:          "MarshalError",
	}

	errorMessages = map[ErrorCode]string{
		DatabaseError:         "An error occurred while reading from the database",
		DatabaseNoRecordError: "No records found",
		UnmarshalError:        "An error occurred while unmarshalling data",
		MarshalError:          "An error occurred while marshaling data",
	}
)

type ErrorResponse struct {
	ErrorReference uuid.UUID `json:"error_reference"`
	Code           ErrorCode `json:"code"`
	ErrorType      string    `json:"error_type"`
	Message        string    `json:"message"`
	Err            error     `json:"-"`
	StackTrace     string    `json:"stack_trace,omitempty"`
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
		errorResponse.StackTrace = fmt.Sprintf("%+v", errors.WithStack(err))
	}

	return errorResponse
}
