package library

import (
	authDomain "github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/mailer"
	orderDomain "github.com/leetatech/leeta_backend/services/order/domain"
	productDomain "github.com/leetatech/leeta_backend/services/product/domain"
	userDomain "github.com/leetatech/leeta_backend/services/user/domain"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Repositories struct {
	OrderRepository   orderDomain.OrderRepository
	UserRepository    userDomain.UserRepository
	AuthRepository    authDomain.AuthRepository
	ProductRepository productDomain.ProductRepository
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
	Domain        string
}

type DefaultErrorResponse struct {
	Data leetError.ErrorResponse `json:"data"`
} // @name DefaultErrorResponse

func GetPaginatedOpts(limit, page int64) *options.FindOptions {
	l := limit
	skip := page*limit - limit
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}

	return &fOpt
}

type BoundryObjectID struct {
	ID string `json:"id" bson:"id"`
}
