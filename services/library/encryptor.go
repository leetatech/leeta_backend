package library

import (
	"errors"
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"unicode"
)

type encryptorHandler struct {
}

var _ EncryptorManager = &encryptorHandler{}

type EncryptorManager interface {
	ComparePasscode(passcode, hashedPasscode string) error
	GenerateFromPasscode(passcode string) ([]byte, error)
	IsValidPassword(s string) error
	IsValidEmailFormat(email string) error
	IsValidHost(email string) error
	IsLeetaDomain(email, leetaDomain string) error
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

		return leetError.ErrorResponseBody(leetError.PasswordValidationError, errors.New("password must contain at least six character long, one uppercase letter, one lowercase letter, one digit, and one special character"))
	default:
		return nil
	}
}

func (e *encryptorHandler) IsValidEmailFormat(email string) error {
	err := checkmail.ValidateFormat(email)
	if err != nil {
		return err
	}

	return nil
}

func (e *encryptorHandler) IsValidHost(email string) error {
	//parts := strings.Split(email, "@")
	//if len(parts) != 2 {
	//	return leetError.ErrorResponseBody(leetError.EmailFormatError, nil)
	//}
	//domain := parts[1]
	err := checkmail.ValidateHost(email)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.ValidEmailHostError, nil)
	}

	return nil
}

func (e *encryptorHandler) IsLeetaDomain(email, leetaDomain string) error {

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return leetError.ErrorResponseBody(leetError.EmailFormatError, nil)
	}
	domain := parts[1]

	fmt.Println("email: ", email)
	fmt.Println("leetaDomain: ", leetaDomain)
	err := checkmail.ValidateHost(email)
	if err != nil {
		fmt.Println(err)
		return leetError.ErrorResponseBody(leetError.ValidEmailHostError, err)
	}

	if strings.ToLower(domain) != strings.ToLower(leetaDomain) {
		return leetError.ErrorResponseBody(leetError.ValidLeetaDomainError, nil)
	}

	return nil
}
