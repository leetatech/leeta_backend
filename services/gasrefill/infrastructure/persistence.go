package infrastructure

import (
	"context"
	"github.com/leetatech/leeta_backend/services/gasrefill/domain"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type refillStoreHandler struct {
	client       *mongo.Client
	databaseName string
	logger       *zap.Logger
}

func (r refillStoreHandler) col(collectionName string) *mongo.Collection {
	return r.client.Database(r.databaseName).Collection(collectionName)
}

func NewRefillPersistence(client *mongo.Client, databaseName string, logger *zap.Logger) domain.GasRefillRepository {
	return &refillStoreHandler{client: client, databaseName: databaseName, logger: logger}
}

func (r refillStoreHandler) RequestRefill(ctx context.Context, request models.GasRefill) error {
	_, err := r.col(models.RefillCollectionName).InsertOne(ctx, request)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}

func (r refillStoreHandler) CreateFees(ctx context.Context, request models.Fees) error {
	_, err := r.col(models.FeesCollectionName).InsertOne(ctx, request)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}

func (r refillStoreHandler) GetFees(ctx context.Context, status models.CartStatuses) (*models.Fees, error) {
	var fee models.Fees
	filter := bson.M{"status": status}
	err := r.col(models.FeesCollectionName).FindOne(ctx, filter).Decode(&fee)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return &fee, nil
}

func (r refillStoreHandler) UpdateFees(ctx context.Context, status models.CartStatuses) error {
	filter := bson.M{"status": models.CartActive}
	update := bson.M{"$set": bson.M{"status": status, "status_ts": time.Now().Unix()}}
	_, err := r.col(models.FeesCollectionName).UpdateMany(ctx, filter, update)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}
