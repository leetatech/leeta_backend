package infrastructure

import (
	"database/sql"
	"fmt"
	"github.com/leetatech/leeta_backend/services/order/domain"

	"go.uber.org/zap"
)

type orderStoreHandler struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewOrderPersistence(db *sql.DB, logger *zap.Logger) domain.OrderRepository {
	return &orderStoreHandler{db: db, logger: logger}
}

func (t orderStoreHandler) CreateOrder(request domain.Order) {
	fmt.Println("Hello")
}
