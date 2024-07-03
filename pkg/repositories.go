package pkg

import (
	"github.com/leetatech/leeta_backend/pkg/config"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/messaging/mailer/awsEmail"
	"github.com/leetatech/leeta_backend/pkg/messaging/mailer/postmarkClient"
	authDomain "github.com/leetatech/leeta_backend/services/auth/domain"
	cartDomain "github.com/leetatech/leeta_backend/services/cart/domain"
	feesDomain "github.com/leetatech/leeta_backend/services/fees/domain"
	orderDomain "github.com/leetatech/leeta_backend/services/order/domain"
	productDomain "github.com/leetatech/leeta_backend/services/product/domain"
	statesDomain "github.com/leetatech/leeta_backend/services/state/domain"
	userDomain "github.com/leetatech/leeta_backend/services/user/domain"
	"go.uber.org/zap"
)

type Repositories struct {
	OrderRepository   orderDomain.OrderRepository
	UserRepository    userDomain.UserRepository
	AuthRepository    authDomain.AuthRepository
	ProductRepository productDomain.ProductRepository
	CartRepository    cartDomain.CartRepository
	FeesRepository    feesDomain.FeesRepository
	StatesRepository  statesDomain.StateRepository
}

type DefaultResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
} // @name DefaultResponse

type DefaultApplicationRequest struct {
	TokenHandler   TokenHandler
	Logger         *zap.Logger
	AllRepository  Repositories
	EmailClient    postmarkClient.MailerClient
	AWSEmailClient awsEmail.AWSEmailClient
	LeetaConfig    config.LeetaConfig
}

type DefaultErrorResponse struct {
	Data leetError.ErrorResponse `json:"data"`
} // @name DefaultErrorResponse
