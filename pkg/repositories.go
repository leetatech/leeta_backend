package pkg

import (
	"github.com/leetatech/leeta_backend/pkg/config"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/pkg/jwtmiddleware"
	"github.com/leetatech/leeta_backend/pkg/mailer/aws"
	authDomain "github.com/leetatech/leeta_backend/services/auth/domain"
	cartDomain "github.com/leetatech/leeta_backend/services/cart/domain"
	feesDomain "github.com/leetatech/leeta_backend/services/fees/domain"
	orderDomain "github.com/leetatech/leeta_backend/services/order/domain"
	productDomain "github.com/leetatech/leeta_backend/services/product/domain"
	statesDomain "github.com/leetatech/leeta_backend/services/state/domain"
	userDomain "github.com/leetatech/leeta_backend/services/user/domain"
)

type RepositoryManager struct {
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

type ApplicationContext struct {
	JwtManager        jwtmiddleware.Manager
	RepositoryManager RepositoryManager
	Domain            string
	MailClient        aws.MailClient
	Config            config.ServerConfig
}

type DefaultErrorResponse struct {
	Data errs.Response `json:"data"`
} // @name DefaultErrorResponse
