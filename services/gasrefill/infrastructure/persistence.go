package infrastructure

import (
	"github.com/leetatech/leeta_backend/services/gasrefill/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type refillStoreHandler struct {
	client       *mongo.Client
	databaseName string
	logger       *zap.Logger
}

func (r *refillStoreHandler) col(collectionName string) *mongo.Collection {
	return r.client.Database(r.databaseName).Collection(collectionName)
}

func NewRefillPersistence(client *mongo.Client, databaseName string, logger *zap.Logger) domain.GasRefillRepository {
	return &refillStoreHandler{client: client, databaseName: databaseName, logger: logger}
}
