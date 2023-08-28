package infrastructure

import (
	"context"
	"github.com/leetatech/leeta_backend/services/library/models"
	"github.com/leetatech/leeta_backend/services/product/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type productStoreHandler struct {
	client       *mongo.Client
	databaseName string
	logger       *zap.Logger
}

func (p productStoreHandler) col(collectionName string) *mongo.Collection {
	return p.client.Database(p.databaseName).Collection(collectionName)
}

func NewProductPersistence(client *mongo.Client, databaseName string, logger *zap.Logger) domain.ProductRepository {
	return &productStoreHandler{client: client, databaseName: databaseName, logger: logger}
}

func (p productStoreHandler) CreateProduct(ctx context.Context, request models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := p.col(models.ProductCollectionName).InsertOne(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func (p productStoreHandler) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	product := &models.Product{}
	filter := bson.M{
		"id": id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := p.col(models.ProductCollectionName).FindOne(ctx, filter).Decode(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p productStoreHandler) GetAllVendorProducts(ctx context.Context, vendorID string) ([]models.Product, error) {
	//TODO implement me
	panic("implement me")
}
