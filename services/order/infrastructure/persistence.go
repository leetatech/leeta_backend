package infrastructure

import (
	"context"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/services/models"
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

func (o orderStoreHandler) CreateOrder(ctx context.Context, request models.Order) error {
	updatedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := o.col(models.OrderCollectionName).InsertOne(updatedCtx, request)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}

func (o orderStoreHandler) UpdateOrderStatus(ctx context.Context, request domain.UpdateOrderStatusRequest) error {
	filter := bson.M{
		"id": request.OrderId,
	}
	update := bson.M{
		"$set": bson.M{
			"status":    request.OrderStatus,
			"status_ts": time.Now().Unix(),
		},
	}

	_, err := o.col(models.OrderCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}
	return nil
}

func (o orderStoreHandler) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	order := &models.Order{}
	filter := bson.M{
		"id": id,
	}

	err := o.col(models.OrderCollectionName).FindOne(ctx, filter).Decode(order)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return order, nil
}

func (o orderStoreHandler) GetCustomerOrdersByStatus(ctx context.Context, request domain.GetCustomerOrders) ([]domain.OrderResponse, error) {
	updatedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{}

	if request.UserId != "" {
		filter["customer_id"] = request.UserId
	}

	if len(request.OrderStatus) > 0 {
		filter["status"] = bson.M{"$in": request.OrderStatus}
	}

	pipeline := makeCustomerPipeline(filter, request.Limit, request.Page)

	cursor, err := o.col(models.OrderCollectionName).Aggregate(updatedCtx, pipeline)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}
	defer cursor.Close(ctx)

	nodes := make([]domain.OrderResponse, cursor.RemainingBatchLength())

	err = cursor.All(ctx, &nodes)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	if len(nodes) == 0 {
		return []domain.OrderResponse{}, nil
	}

	return nodes, err
}

func makeCustomerPipeline(filter bson.M, limit, page int64) []bson.M {

	skip := page*limit - limit
	customerFilter := bson.M{"$expr": bson.M{"$eq": []string{"$id", "$$customer_id"}}}
	customerPipeline := []bson.M{
		{
			"$match": customerFilter,
		},
	}

	productFilter := bson.M{"$expr": bson.M{"$eq": []string{"$id", "$$product_id"}}}
	productPipeline := []bson.M{
		{
			"$match": productFilter,
		},
	}

	return []bson.M{
		{
			"$match": filter,
		},
		{
			"$lookup": bson.M{
				"from":         models.UsersCollectionName,
				"let":          bson.M{"customer_id": "$customer_id"},
				"localField":   "customer_id",
				"foreignField": "id",
				"pipeline":     customerPipeline,
				"as":           "customer",
			},
		},
		{
			"$unwind": "$customer",
		},
		{
			"$lookup": bson.M{
				"from":         models.ProductCollectionName,
				"let":          bson.M{"product_id": "$product_id"},
				"localField":   "product_id",
				"foreignField": "id",
				"pipeline":     productPipeline,
				"as":           "products",
			},
		},
		{
			"$unwind": "$products",
		},
		{
			"$skip": skip,
		},
		{
			"$limit": limit,
		},
	}
}
