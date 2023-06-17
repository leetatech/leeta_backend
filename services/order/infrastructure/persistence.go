package infrastructure

import (
	"context"
	"fmt"
	"github.com/leetatech/leeta_backend/services/order/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"go.uber.org/zap"
)

type orderStoreHandler struct {
	client       *mongo.Client
	databaseName string
	logger       *zap.Logger
}

func (o orderStoreHandler) col(collectionName string) *mongo.Collection {
	return o.client.Database(o.databaseName).Collection(collectionName)
}

func NewOrderPersistence(client *mongo.Client, databaseName string, logger *zap.Logger) domain.OrderRepository {
	return &orderStoreHandler{client: client, databaseName: databaseName, logger: logger}
}

func (o orderStoreHandler) CreateOrder(request domain.Order) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	type f struct {
		Name string `json:"name"`
	}
	_, err := o.col("leet").InsertOne(ctx, f{Name: "Tim"})
	if err != nil {
		fmt.Println(err)
		fmt.Println("col:error")
	}
	fmt.Println("yes!")
}
