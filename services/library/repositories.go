package library

import (
	authDomain "github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/mailer"
	orderDomain "github.com/leetatech/leeta_backend/services/order/domain"
	userDomain "github.com/leetatech/leeta_backend/services/user/domain"
	"go.uber.org/zap"
)

type Repositories struct {
	OrderRepository orderDomain.OrderRepository
	UserRepository  userDomain.UserRepository
	AuthRepository  authDomain.AuthRepository
}

type DefaultResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
} // @name DefaultResponse

type DefaultApplicationRequest struct {
	TokenHandler  TokenHandler
	Logger        *zap.Logger
	AllRepository Repositories
	EmailClient   mailer.MailerClient
}

type DefaultErrorResponse struct {
	Data leetError.ErrorResponse `json:"data"`
} // @name DefaultErrorResponse
