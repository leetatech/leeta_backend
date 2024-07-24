package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/leetatech/leeta_backend/pkg/database"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/order/domain"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type orderStoreHandler struct {
	client       *mongo.Client
	databaseName string
}

func (o *orderStoreHandler) col(collectionName string) *mongo.Collection {
	return o.client.Database(o.databaseName).Collection(collectionName)
}

func New(client *mongo.Client, databaseName string) domain.OrderRepository {
	return &orderStoreHandler{client: client, databaseName: databaseName}
}

func (o *orderStoreHandler) Create(ctx context.Context, request models.Order) error {
	updatedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := o.col(models.OrderCollectionName).InsertOne(updatedCtx, request)
	if err != nil {
		return errs.Body(errs.DatabaseError, err)
	}

	return nil
}

func (o *orderStoreHandler) UpdateStatus(ctx context.Context, request domain.PersistOrderUpdate) error {
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
		return errs.Body(errs.DatabaseError, err)
	}
	return nil
}

func (o *orderStoreHandler) OrderByID(ctx context.Context, id string) (*models.Order, error) {
	order := &models.Order{}
	filter := bson.M{
		"id": id,
	}

	err := o.col(models.OrderCollectionName).FindOne(ctx, filter).Decode(order)
	if err != nil {
		return nil, errs.Body(errs.DatabaseError, err)
	}

	return order, nil
}

func (o *orderStoreHandler) OrdersByStatus(ctx context.Context, request domain.GetCustomerOrders) ([]domain.Response, error) {
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
		return nil, errs.Body(errs.DatabaseError, err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Debug().Msgf("error closing mongo cursor %v", err)
		}
	}(cursor, ctx)

	nodes := make([]domain.Response, cursor.RemainingBatchLength())

	err = cursor.All(ctx, &nodes)
	if err != nil {
		return nil, errs.Body(errs.DatabaseError, err)
	}

	if len(nodes) == 0 {
		return []domain.Response{}, nil
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

func (o *orderStoreHandler) Orders(ctx context.Context, request query.ResultSelector, userId string) (orders []models.Order, totalResults uint64, err error) {
	updatedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{"customer_id": userId}
	if request.Filter != nil {
		userFilter := database.BuildMongoFilterQuery(request.Filter, nil)
		for key, value := range userFilter {
			filter[key] = value
		}
	}

	totalRecord, err := o.col(models.OrderCollectionName).CountDocuments(updatedCtx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Calculate skip value and ensure it's non-negative
	skip := int64(request.Paging.PageSize * request.Paging.PageIndex)
	if skip < 0 {
		skip = 0
	}

	limit := int64(request.Paging.PageSize)
	pagingOptions := options.Find().SetSkip(skip).SetLimit(limit)

	extraDocumentCursor, err := o.col(models.OrderCollectionName).Find(updatedCtx, filter, options.Find().SetSkip(skip+limit).SetLimit(1))
	if err != nil {
		return nil, 0, fmt.Errorf("error finding orders: %w", err)
	}
	defer func(extraDocumentCursor *mongo.Cursor, ctx context.Context) {
		err = extraDocumentCursor.Close(ctx)
		if err != nil {
			log.Debug().Msgf("error closing mongo cursor %v", err)
		}
	}(extraDocumentCursor, ctx)

	cursor, err := o.col(models.OrderCollectionName).Find(updatedCtx, filter, pagingOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("error finding orders in mongo collection: %w", err)
	}
	orders = make([]models.Order, cursor.RemainingBatchLength())
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, 0, fmt.Errorf("error getting orders: %w", err)
	}

	return orders, uint64(totalRecord), nil
}

func (o *orderStoreHandler) OrderStatusHistory(ctx context.Context, orderId string) ([]models.StatusHistory, error) {
	updatedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{
		"id": orderId,
	}

	order := &models.Order{}
	err := o.col(models.OrderCollectionName).FindOne(updatedCtx, filter).Decode(order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errs.Body(errs.DatabaseNoRecordError, fmt.Errorf("order with id %s not found", orderId))
		}
		return nil, errs.Body(errs.DatabaseError, err)
	}

	return order.StatusHistory, nil
}
