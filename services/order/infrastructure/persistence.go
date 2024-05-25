package infrastructure

import (
	"context"
	"github.com/leetatech/leeta_backend/pkg/database"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/query"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/order/domain"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (o orderStoreHandler) UpdateOrderStatus(ctx context.Context, request domain.PersistOrderUpdate) error {
	filter := bson.M{
		"id": request.OrderId,
	}
	update := bson.M{
		"$set": bson.M{
			"status":    request.OrderStatus,
			"reason":    request.Reason,
			"status_ts": time.Now().Unix(),
		},
		"$push": bson.M{
			"status_history": request.StatusHistory,
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
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Debug().Msgf("error closing mongo cursor %v", err)
		}
	}(cursor, ctx)

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

func (o orderStoreHandler) ListOrders(ctx context.Context, request query.ResultSelector, userId string) (orders []models.Order, totalResults uint64, err error) {
	updatedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var filter bson.M

	var pagingOptions *options.FindOptions
	if request.Filter != nil {
		filter = database.BuildMongoFilterQuery(request.Filter)
		filter["customer_id"] = userId
	}

	totalRecord, err := o.col(models.OrderCollectionName).CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	pagingOptions = database.GetPaginatedOpts(int64(request.Paging.PageSize), int64(request.Paging.PageIndex))

	extraDocumentCursor, err := o.col(models.OrderCollectionName).Find(updatedCtx, filter, options.Find().SetSkip(*pagingOptions.Skip+*pagingOptions.Limit).SetLimit(1))
	if err != nil {
		o.logger.Error("error getting extra document", zap.Error(err))
		return nil, 0, err
	}
	defer func(extraDocumentCursor *mongo.Cursor, ctx context.Context) {
		err = extraDocumentCursor.Close(ctx)
		if err != nil {
			log.Debug().Msgf("error closing mongo cursor %v", err)
		}
	}(extraDocumentCursor, ctx)

	cursor, err := o.col(models.OrderCollectionName).Find(updatedCtx, filter, pagingOptions)
	if err != nil {
		o.logger.Error("error getting orders", zap.Error(err))
		return nil, 0, err
	}
	orders = make([]models.Order, cursor.RemainingBatchLength())
	if err = cursor.All(ctx, &orders); err != nil {
		o.logger.Error("error getting orders", zap.Error(err))
		return nil, 0, err
	}

	return orders, uint64(totalRecord), nil
}

func (o orderStoreHandler) ListOrderStatusHistory(ctx context.Context, orderId string) ([]models.StatusHistory, error) {
	updatedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{
		"id": orderId,
	}

	order := &models.Order{}
	err := o.col(models.OrderCollectionName).FindOne(updatedCtx, filter).Decode(order)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return order.StatusHistory, nil
}
