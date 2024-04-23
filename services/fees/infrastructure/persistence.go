package infrastructure

import (
	"context"
	"github.com/leetatech/leeta_backend/pkg/database"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/query"
	"github.com/leetatech/leeta_backend/services/fees/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (f feeStoreHandler) CreateFees(ctx context.Context, request models.Fee) error {
	_, err := f.col(models.FeesCollectionName).InsertOne(ctx, request)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}

func (f feeStoreHandler) GetFeesByStatus(ctx context.Context, status models.FeesStatuses) ([]models.Fee, error) {
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

func (f feeStoreHandler) UpdateFees(ctx context.Context, status models.FeesStatuses, feeType models.FeeType, lga models.LGA, productID string) error {
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
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}

func (f feeStoreHandler) GetFeeByProductID(ctx context.Context, productID string, status models.FeesStatuses) (*models.Fee, error) {
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

func (f feeStoreHandler) GetTypedFees(ctx context.Context, request query.ResultSelector) ([]models.Fee, uint64, error) {
	opt := database.GetPaginatedOpts(int64(request.Paging.PageSize), int64(request.Paging.PageIndex))
	var filterQuery bson.M
	if request.Filter != nil {
		filterQuery = database.BuildMongoFilterQuery(request.Filter)
	}
	totalRecord, err := f.col(models.FeesCollectionName).CountDocuments(ctx, filterQuery)
	if err != nil {
		return nil, 0, err
	}

	extraDocumentCursor, err := f.col(models.FeesCollectionName).Find(ctx, filterQuery, options.Find().SetSkip(*opt.Skip+*opt.Limit).SetLimit(1))
	if err != nil {
		f.logger.Error("error getting extra document", zap.Error(err))
		return nil, 0, err
	}
	defer func(extraDocumentCursor *mongo.Cursor, ctx context.Context) {
		err = extraDocumentCursor.Close(ctx)
		if err != nil {
			log.Debug().Msgf("error closing mongo cursor %v", err)
		}
	}(extraDocumentCursor, ctx)

	cursor, err := f.col(models.FeesCollectionName).Find(ctx, filterQuery, opt)
	if err != nil {
		f.logger.Error("error getting fees", zap.Error(err))
		return nil, 0, err
	}

	fees := make([]models.Fee, cursor.RemainingBatchLength())
	if err := cursor.All(ctx, &fees); err != nil {
		f.logger.Error("error getting fees", zap.Error(err))
		return nil, 0, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Debug().Msgf("error closing mongo cursur %v", err)
		}
	}(cursor, ctx)

	return fees, uint64(totalRecord), nil
}
