package infrastructure

import (
	"context"
	"github.com/leetatech/leeta_backend/pkg/database"
	"github.com/leetatech/leeta_backend/pkg/filter"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/services/cart/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
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

func (c *CartStoreHandler) DeleteCartItem(ctx context.Context, cartItemID string, itemTotalCost float64) error {
	filter := bson.M{"cart_items.id": cartItemID}

	update := bson.M{
		"$pull": bson.M{
			"cart_items": bson.M{"id": cartItemID},
		},
		"$inc": bson.M{"total": -itemTotalCost},
		"$set": bson.M{"status_ts": time.Now().Unix()},
	}

	_, err := c.col(models.CartsCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (c *CartStoreHandler) DeleteCart(ctx context.Context, id string) error {
	filter := bson.M{"id": id}

	_, err := c.col(models.CartsCollectionName).DeleteOne(ctx, filter)
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

func (c *CartStoreHandler) GetCartByCartItemID(ctx context.Context, cartItemID string) (*models.Cart, error) {
	var cart models.Cart
	filter := bson.M{"cart_items.id": cartItemID}

	err := c.col(models.CartsCollectionName).FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (c *CartStoreHandler) GetPaginatedCart(ctx context.Context, request filter.ResultSelector, userID string) (*domain.ListCartResponse, error) {
	opt := database.GetPaginatedOpts(int64(request.Paging.PageSize), int64(request.Paging.PageIndex))

	query := database.BuildMongoFilterQuery(request.Filter)
	query["customer_id"] = userID

	pipeline := mongo.Pipeline{
		{
			{"$match", query},
		},
		{
			{"$project", bson.M{
				"id":    1,
				"total": 1,
				"total_records": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$isArray": "$cart_items"},
						"then": bson.M{"$size": "$cart_items"},
						"else": 0,
					},
				},
				"cart_items": bson.M{"$slice": []interface{}{"$cart_items", opt.Skip, opt.Limit}},
			}},
		},
	}

	cursor, err := c.col(models.CartsCollectionName).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var cartResponse domain.CartResponse
	if cursor.Next(ctx) {
		if err = cursor.Decode(&cartResponse); err != nil {
			return nil, err
		}
	}

	return &domain.ListCartResponse{
		Cart:        cartResponse,
		HasNextPage: (request.Paging.PageIndex * request.Paging.PageSize) < cartResponse.TotalRecords,
	}, nil
}
