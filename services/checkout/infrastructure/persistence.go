package infrastructure

import (
	"context"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/services/checkout/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type checkoutStoreHandler struct {
	client       *mongo.Client
	databaseName string
	logger       *zap.Logger
}

func (c checkoutStoreHandler) col(collectionName string) *mongo.Collection {
	return c.client.Database(c.databaseName).Collection(collectionName)
}

func NewCheckoutPersistence(client *mongo.Client, databaseName string, logger *zap.Logger) domain.CheckoutRepository {
	return &checkoutStoreHandler{client: client, databaseName: databaseName, logger: logger}
}

func (c checkoutStoreHandler) RequestCheckout(ctx context.Context, request models.Checkout) error {
	_, err := c.col(models.CheckoutCollectionName).InsertOne(ctx, request)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}
