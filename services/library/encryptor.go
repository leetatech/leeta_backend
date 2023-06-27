package library

import (
	"errors"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

type encryptorHandler struct {
}

var _ EncryptorManager = &encryptorHandler{}

type EncryptorManager interface {
	ComparePasscode(passcode, hashedPasscode string) error
	GenerateFromPasscode(passcode string) ([]byte, error)
	IsValidPassword(s string) error
}

func NewEncryptor() EncryptorManager {
	return &encryptorHandler{}
}

func (e *encryptorHandler) GenerateFromPasscode(passcode string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(passcode), bcrypt.DefaultCost)
}

func (e *encryptorHandler) ComparePasscode(passcode, hashedPasscode string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPasscode), []byte(passcode))
}

func (e *encryptorHandler) IsValidPassword(s string) error {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	switch {
	case len(s) >= 6:
		hasMinLen = true
	default:
		hasMinLen = false
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	switch {
	case !hasMinLen, !hasUpper, !hasLower, !hasNumber, !hasSpecial:

		return leetError.ErrorResponseBody(leetError.DatabaseNoRecordError, errors.New("password must contain at least six character long, one uppercase letter, one lowercase letter, one digit, and one special character"))
	default:
		return nil
	}
}
