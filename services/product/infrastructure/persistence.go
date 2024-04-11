package infrastructure

import (
	"context"
	"errors"
	"github.com/leetatech/leeta_backend/pkg/database"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/query"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/product/domain"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	updatedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := p.col(models.ProductCollectionName).InsertOne(updatedCtx, request)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}

func (p productStoreHandler) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	product := &models.Product{}
	filter := bson.M{
		"id": id,
	}

	updatedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := p.col(models.ProductCollectionName).FindOne(updatedCtx, filter).Decode(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p productStoreHandler) GetAllVendorProducts(ctx context.Context, request domain.GetVendorProductsRequest) (*domain.GetVendorProductsResponse, error) {
	filter := bson.M{}
	filter["vendor_id"] = request.VendorID
	if len(request.ProductStatus) > 0 {
		filter["status"] = bson.M{"$in": request.ProductStatus}
	}

	opts := database.GetPaginatedOpts(request.Limit, request.Page)

	updatedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cursor, err := p.col(models.ProductCollectionName).Find(updatedCtx, filter, opts)
	if err != nil {
		return nil, err
	}
	products := make([]models.Product, cursor.RemainingBatchLength())

	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	var hasNextPage bool

	// TODO get the remaining batch the right way
	//if cursor.RemainingBatchLength() > 0 {
	//	hasNextPage = true
	//}

	return &domain.GetVendorProductsResponse{
		Products:    products,
		HasNextPage: hasNextPage,
	}, nil
}

func (p productStoreHandler) ListProducts(ctx context.Context, request *query.ResultSelector) (*query.ResponseListWithMetadata[models.Product], error) {
	updatedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var filter bson.M
	var pagingOptions *options.FindOptions
	if request == nil {
		return nil, errors.New("result selector request is required")
	}
	if request.Filter != nil {
		filter = database.BuildMongoFilterQuery(request.Filter)
	}

	totalRecord, err := p.col(models.ProductCollectionName).CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	pagingOptions = database.GetPaginatedOpts(int64(request.Paging.PageSize), int64(request.Paging.PageIndex))

	extraDocumentCursor, err := p.col(models.ProductCollectionName).Find(updatedCtx, filter, options.Find().SetSkip(*pagingOptions.Skip+*pagingOptions.Limit).SetLimit(1))
	if err != nil {
		p.logger.Error("error getting extra document", zap.Error(err))
		return nil, err
	}
	defer func(extraDocumentCursor *mongo.Cursor, ctx context.Context) {
		err = extraDocumentCursor.Close(ctx)
		if err != nil {
			log.Debug().Msgf("error closing mongo cursor %v", err)
		}
	}(extraDocumentCursor, ctx)
	hasNextPage := extraDocumentCursor.Next(ctx)

	cursor, err := p.col(models.ProductCollectionName).Find(updatedCtx, filter, pagingOptions)
	if err != nil {
		p.logger.Error("error getting products", zap.Error(err))
		return nil, err
	}
	products := make([]models.Product, cursor.RemainingBatchLength())
	if err = cursor.All(ctx, &products); err != nil {
		p.logger.Error("error getting products", zap.Error(err))
		return nil, err
	}

	return &query.ResponseListWithMetadata[models.Product]{
		Metadata: query.NewMetadata(*request, uint64(totalRecord), hasNextPage),
		Data:     products,
	}, nil
}
