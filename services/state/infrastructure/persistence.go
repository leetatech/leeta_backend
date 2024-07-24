package infrastructure

import (
	"context"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/state/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type StateStoreHandler struct {
	client       *mongo.Client
	databaseName string
	logger       *zap.Logger
}

func New(client *mongo.Client, databaseName string) domain.StateRepository {
	return &StateStoreHandler{client: client, databaseName: databaseName}
}

func (s *StateStoreHandler) SaveStates(ctx context.Context, states []any) error {
	_, err := s.col(models.NGNStatesCollectionName).InsertMany(ctx, states)
	if err != nil {
		return errs.Body(errs.DatabaseError, err)
	}

	return nil
}

func (s *StateStoreHandler) GetState(ctx context.Context, name string) (models.State, error) {
	var state models.State
	filter := bson.M{"name": name}

	err := s.col(models.NGNStatesCollectionName).FindOne(ctx, filter).Decode(&state)
	if err != nil {
		return state, err
	}

	return state, nil
}

func (s *StateStoreHandler) GetAllStates(ctx context.Context) ([]models.State, error) {
	cursor, err := s.col(models.NGNStatesCollectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	states := make([]models.State, cursor.RemainingBatchLength())
	if err = cursor.All(ctx, &states); err != nil {
		return nil, err
	}

	return states, nil
}

func (s *StateStoreHandler) col(collectionName string) *mongo.Collection {
	return s.client.Database(s.databaseName).Collection(collectionName)
}
