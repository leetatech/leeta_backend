package infrastructure

import (
	"context"
	"github.com/leetatech/leeta_backend/services/gasrefill/domain"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
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
