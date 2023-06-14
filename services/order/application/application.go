package application

import (
	"context"
	"fmt"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/order/domain"
)

/**




// PLACEHOLDERS




**/

var (
	universalPassword          string = "leeta"
	universalEncryptedPassword string
)

type orderAppHandler struct {
	tokenHandler  library.TokenHandler
	encryptor     library.EncryptorManager
	allRepository library.Repositories
}

type OrderApplication interface {
	CreateOrder(ctx context.Context, request domain.Order) (*library.DefaultResponse, error)
}

func NewOrderApplication(tokenHandler library.TokenHandler, allRepository library.Repositories) OrderApplication {
	encryptor := library.NewEncryptor()

	// random password
	encryptPassword(encryptor)

	return &orderAppHandler{
		tokenHandler:  tokenHandler,
		encryptor:     encryptor,
		allRepository: allRepository,
	}
}

func (t orderAppHandler) CreateOrder(ctx context.Context, request domain.Order) (*library.DefaultResponse, error) {
	claims, err := t.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		fmt.Println("unable to get claims")
	}
	if t.validatePin(request.Status) != nil {
		return nil, fmt.Errorf("%s", "invalid pin credential")
	}
	fmt.Println(claims)
	return nil, err
}

func encryptPassword(encryptor library.EncryptorManager) {
	byteValueOfPassword, _ := encryptor.GenerateFromPasscode(universalPassword)
	universalEncryptedPassword = string(byteValueOfPassword)
}

func (t orderAppHandler) validatePin(pin string) error {
	return t.encryptor.ComparePasscode(pin, "")
}
