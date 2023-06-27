package library

import (
	authDomain "github.com/leetatech/leeta_backend/services/auth/domain"
	orderDomain "github.com/leetatech/leeta_backend/services/order/domain"
	userDomain "github.com/leetatech/leeta_backend/services/user/domain"
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
