package infrastructure

import (
	"context"
	"github.com/leetatech/leeta_backend/services/fees/domain"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type feeStoreHandler struct {
	client       *mongo.Client
	databaseName string
	logger       *zap.Logger
}

func (f feeStoreHandler) col(collectionName string) *mongo.Collection {
	return f.client.Database(f.databaseName).Collection(collectionName)
}

func NewFeesPersistence(client *mongo.Client, databaseName string, logger *zap.Logger) domain.FeesRepository {
	return &feeStoreHandler{client: client, databaseName: databaseName, logger: logger}
}

func (f feeStoreHandler) CreateFees(ctx context.Context, request models.Fees) error {
	_, err := f.col(models.FeesCollectionName).InsertOne(ctx, request)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}

func (f feeStoreHandler) GetFees(ctx context.Context, status models.FeesStatuses) ([]models.Fees, error) {
	filter := bson.M{"status": status}

	cursor, err := f.col(models.FeesCollectionName).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	fees := make([]models.Fees, cursor.RemainingBatchLength())
	if err := cursor.All(ctx, &fees); err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	return fees, nil
}

func (f feeStoreHandler) UpdateFees(ctx context.Context, status models.FeesStatuses) error {
	filter := bson.M{"status": models.CartActive}
	update := bson.M{"$set": bson.M{"status": status, "status_ts": time.Now().Unix()}}
	_, err := f.col(models.FeesCollectionName).UpdateMany(ctx, filter, update)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}

func (f feeStoreHandler) GetFeeByProductID(ctx context.Context, productID string, status models.FeesStatuses) (*models.Fees, error) {
	filter := bson.M{"product_id": productID, "status": status}
	fee := &models.Fees{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := f.col(models.FeesCollectionName).FindOne(ctx, filter).Decode(fee)
	if err != nil {
		return nil, err
	}

	return fee, nil
}
