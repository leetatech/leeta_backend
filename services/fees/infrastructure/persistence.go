package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/leetatech/leeta_backend/pkg/database"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/services/fees/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type feeStoreHandler struct {
	client       *mongo.Client
	databaseName string
}

func (f *feeStoreHandler) col(collectionName string) *mongo.Collection {
	return f.client.Database(f.databaseName).Collection(collectionName)
}

func New(client *mongo.Client, databaseName string) domain.FeesRepository {
	return &feeStoreHandler{client: client, databaseName: databaseName}
}

func (f *feeStoreHandler) Create(ctx context.Context, request models.Fee) error {
	_, err := f.col(models.FeesCollectionName).InsertOne(ctx, request)
	if err != nil {
		return errs.Body(errs.DatabaseError, err)
	}

	return nil
}

func (f *feeStoreHandler) FeesByStatus(ctx context.Context, status models.FeesStatuses) ([]models.Fee, error) {
	filter := bson.M{"status": status}

	cursor, err := f.col(models.FeesCollectionName).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	fees := make([]models.Fee, cursor.RemainingBatchLength())
	if err := cursor.All(ctx, &fees); err != nil {
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Debug().Msgf("error closing mongo cursur %v", err)
		}
	}(cursor, ctx)

	return fees, nil
}

func (f *feeStoreHandler) Update(ctx context.Context, status models.FeesStatuses, feeType models.FeeType, lga models.LGA, productID string) error {
	filter := bson.M{}
	if status != "" {
		filter["status"] = models.FeesActive
	}

	if feeType != "" {
		filter["fee_type"] = feeType
	}

	if lga != (models.LGA{}) {
		filter["lga"] = lga
	}

	if productID != "" {
		filter["product_id"] = productID
	}

	update := bson.M{"$set": bson.M{"status": status, "status_ts": time.Now().Unix()}}
	_, err := f.col(models.FeesCollectionName).UpdateMany(ctx, filter, update)
	if err != nil {
		return errs.Body(errs.DatabaseError, err)
	}

	return nil
}

func (f *feeStoreHandler) ByProductID(ctx context.Context, productID string, status models.FeesStatuses) (*models.Fee, error) {
	filter := bson.M{"product_id": productID, "status": status}
	fee := &models.Fee{}

	newCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := f.col(models.FeesCollectionName).FindOne(newCtx, filter).Decode(fee)
	if err != nil {
		return nil, err
	}

	return fee, nil
}

func (f *feeStoreHandler) Fees(ctx context.Context, request query.ResultSelector) ([]models.Fee, uint64, error) {
	feesFilterMapping := map[string]string{
		"lga": "lga.lga",
	}
	var filterQuery bson.M
	if request.Filter != nil {
		filterQuery = database.BuildMongoFilterQuery(request.Filter, feesFilterMapping)
	}

	cursor, err := f.col(models.FeesCollectionName).Find(ctx, filterQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting fees: %w", err)
	}

	fees := make([]models.Fee, cursor.RemainingBatchLength())
	if err := cursor.All(ctx, &fees); err != nil {
		return nil, 0, fmt.Errorf("error getting remaining batch length of fees: %w", err)
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Debug().Msgf("error closing mongo cursur %v", err)
		}
	}(cursor, ctx)

	return fees, uint64(len(fees)), nil
}
