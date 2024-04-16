package infrastructure

import (
	"context"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/services/address/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type AddressStoreHandler struct {
	client       *mongo.Client
	databaseName string
	logger       *zap.Logger
}

func (a *AddressStoreHandler) Upsert(ctx context.Context, state models.State) error {
	filter := bson.M{
		"name": state.Name,
	}
	update := bson.M{
		"$setOnInsert": state,
	}
	_, err := a.col(models.AddressCollectionName).UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}

func (a *AddressStoreHandler) Update(ctx context.Context, state models.State) error {
	_, err := a.col(models.AddressCollectionName).UpdateOne(ctx, bson.M{"id": state.Name}, bson.M{"$set": state})
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}
	return nil
}

func (a *AddressStoreHandler) GetState(ctx context.Context, name string) (models.State, error) {
	var state models.State
	filter := bson.M{"name": name}

	err := a.col(models.AddressCollectionName).FindOne(ctx, filter).Decode(&state)
	if err != nil {
		return state, err
	}

	return state, nil
}

func (a *AddressStoreHandler) GetAllStates(ctx context.Context) ([]models.State, error) {
	cursor, err := a.col(models.AddressCollectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	states := make([]models.State, cursor.RemainingBatchLength())
	if err = cursor.All(ctx, &states); err != nil {
		return nil, err
	}

	return states, nil
}

func (a *AddressStoreHandler) col(collectionName string) *mongo.Collection {
	return a.client.Database(a.databaseName).Collection(collectionName)
}

func NewAddressPersistence(client *mongo.Client, databaseName string, logger *zap.Logger) domain.AddressRepository {
	return &AddressStoreHandler{client: client, databaseName: databaseName, logger: logger}
}
