package infrastructure

import (
	"context"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/models"
	"github.com/leetatech/leeta_backend/services/order/domain"
	"go.mongodb.org/mongo-driver/bson"
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

func (o orderStoreHandler) CreateOrder(request domain.OrderRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := o.col(models.OrderCollectionName).InsertOne(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (o orderStoreHandler) UpdateOrderStatus(request domain.UpdateOrderStatusRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"id": request.OrderId,
	}
	update := bson.M{
		"$set": bson.M{
			"status": request.OrderStatus,
		},
	}

	_, err := o.col(models.OrderCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}
	return nil
}

func (o orderStoreHandler) GetOrderByID(id string) (*models.Order, error) {
	order := &models.Order{}
	filter := bson.M{
		"id": id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := o.col(models.OrderCollectionName).FindOne(ctx, filter).Decode(order)
	if err != nil {
		return nil, err
	}

	return order, nil

}
