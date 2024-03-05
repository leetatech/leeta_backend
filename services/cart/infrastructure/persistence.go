package infrastructure

import (
	"context"
	"github.com/leetatech/leeta_backend/services/cart/domain"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type CartStoreHandler struct {
	client       *mongo.Client
	databaseName string
	logger       *zap.Logger
}

func (c *CartStoreHandler) col(collectionName string) *mongo.Collection {
	return c.client.Database(c.databaseName).Collection(collectionName)
}

func NewCartPersistence(client *mongo.Client, databaseName string, logger *zap.Logger) domain.CartRepository {
	return &CartStoreHandler{client: client, databaseName: databaseName, logger: logger}
}

func (c *CartStoreHandler) AddToCart(ctx context.Context, request models.Cart) error {
	_, err := c.col(models.CartsCollectionName).InsertOne(ctx, request)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}

func (c *CartStoreHandler) GetCartByCustomerID(ctx context.Context, customerID string) (*models.Cart, error) {
	var cart models.Cart
	filter := bson.M{"customer_id": customerID, "status": models.CartActive}

	err := c.col(models.CartsCollectionName).FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (c *CartStoreHandler) UpdateCart(ctx context.Context, request models.Cart) error {
	_, err := c.col(models.CartsCollectionName).UpdateOne(ctx, bson.M{"id": request.ID}, bson.M{"$set": request})
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}
	return nil
}

func (c *CartStoreHandler) AddToCartItem(ctx context.Context, cartID string, cartItems models.CartItem, total float64, statusTs int64) error {
	filter := bson.M{"id": cartID}
	update := bson.M{"$push": bson.M{"cart_items": cartItems}, "$set": bson.M{"total": total, "status_ts": statusTs}}

	_, err := c.col(models.CartsCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (c *CartStoreHandler) DeleteCartItem(ctx context.Context, cartID, cartItemID string) error {
	filter := bson.M{"id": cartID}

	update := bson.M{"$pull": bson.M{"cart_items": bson.M{"id": cartItemID}}}

	_, err := c.col(models.CartsCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (c *CartStoreHandler) InactivateCart(ctx context.Context, id string) error {
	filter := bson.M{"id": id}

	update := bson.M{"$set": bson.M{"status": models.CartInactive}}

	_, err := c.col(models.CartsCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (c *CartStoreHandler) GetCartByDeviceID(ctx context.Context, deviceID string) (*models.Cart, error) {
	var cart models.Cart
	filter := bson.M{"device_id": deviceID, "status": models.CartActive}

	err := c.col(models.CartsCollectionName).FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		return nil, err
	}

	return &cart, nil
}
