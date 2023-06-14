package library

import (
	"golang.org/x/crypto/bcrypt"
)

type encryptorHandler struct {
}

var _ EncryptorManager = &encryptorHandler{}

type EncryptorManager interface {
	ComparePasscode(passcode, hashedPasscode string) error
	GenerateFromPasscode(passcode string) ([]byte, error)
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
