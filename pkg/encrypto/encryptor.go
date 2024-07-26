package encrypto

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"unicode"
)

type encryptorHandler struct {
}

var _ Manager = &encryptorHandler{}

type Manager interface {
	ComparePasscode(passcode, hashedPasscode string) error
	Generate(passcode string) ([]byte, error)
	ValidatePasswordStrength(s string) error
	ValidateEmailFormat(email string) error
	ValidateDomain(email, leetaDomain string) error
}

func New() Manager {
	return &encryptorHandler{}
}

func (e *encryptorHandler) Generate(passcode string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(passcode), bcrypt.DefaultCost)
}

func (e *encryptorHandler) ComparePasscode(passcode, hashedPasscode string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPasscode), []byte(passcode))
}

func (e *encryptorHandler) ValidatePasswordStrength(password string) error {
	const minLen = 6
	var hasUpper, hasLower, hasNumber, hasSpecial bool

	if len(password) < minLen {
		return errs.Body(errs.PasswordValidationError, errors.New("password must be at least six characters long"))
	}

	for _, char := range password {
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

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errs.Body(errs.PasswordValidationError, errors.New("password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character"))
	}

	return nil
}

func (e *encryptorHandler) ValidateEmailFormat(email string) error {
	if err := checkmail.ValidateFormat(email); err != nil {
		return err
	}
	return nil
}

func (e *encryptorHandler) ValidateDomain(email, domainString string) error {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errs.Body(errs.EmailFormatError, nil)
	}
	domain := parts[1]

	if err := checkmail.ValidateHost(email); err != nil {
		return errs.Body(errs.ValidEmailHostError, err)
	}

	if strings.EqualFold(domain, domainString) {
		return errs.Body(errs.ValidLeetaDomainError, nil)
	}

	return nil
}
