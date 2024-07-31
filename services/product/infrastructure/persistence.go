package infrastructure

import (
	"context"
	"fmt"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/leetatech/leeta_backend/pkg/database"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/product/domain"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type productStoreHandler struct {
	client       *mongo.Client
	databaseName string
}

func (p productStoreHandler) col(collectionName string) *mongo.Collection {
	return p.client.Database(p.databaseName).Collection(collectionName)
}

func New(client *mongo.Client, databaseName string) domain.ProductRepository {
	return &productStoreHandler{client: client, databaseName: databaseName}
}

func (p productStoreHandler) Create(ctx context.Context, request models.Product) error {
	updatedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := p.col(models.ProductCollectionName).InsertOne(updatedCtx, request)
	if err != nil {
		return errs.Body(errs.DatabaseError, err)
	}

	return nil
}

func (p productStoreHandler) Product(ctx context.Context, id string) (models.Product, error) {
	filter := bson.M{
		"id": id,
	}

	product := models.Product{}

	updatedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := p.col(models.ProductCollectionName).FindOne(updatedCtx, filter).Decode(&product)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (p productStoreHandler) VendorProducts(ctx context.Context, request domain.GetVendorProductsRequest) ([]models.Product, error) {
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

	// TODO get the remaining batch the right way
	//if cursor.RemainingBatchLength() > 0 {
	//	hasNextPage = true
	//}

	return products, nil
}

func (p productStoreHandler) ListProducts(ctx context.Context, request query.ResultSelector) (products []models.Product, totalResults uint64, err error) {
	updatedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var filter bson.M
	var pagingOptions *options.FindOptions
	if request.Filter != nil {
		filter = database.BuildMongoFilterQuery(request.Filter, nil)
	}

	totalRecord, err := p.col(models.ProductCollectionName).CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// TODO: uncomment code below when bug is fixed
	//pagingOptions = database.GetPaginatedOpts(int64(request.Paging.PageSize), int64(request.Paging.PageIndex))

	//extraDocumentCursor, err := p.col(models.ProductCollectionName).Find(updatedCtx, filter, options.Find().SetSkip(*pagingOptions.Skip+*pagingOptions.Limit).SetLimit(1))
	extraDocumentCursor, err := p.col(models.ProductCollectionName).Find(updatedCtx, filter, options.Find().SetLimit(1))
	if err != nil {
		return nil, 0, fmt.Errorf("error getting extra documents: %w", err)
	}
	defer func(extraDocumentCursor *mongo.Cursor, ctx context.Context) {
		err = extraDocumentCursor.Close(ctx)
		if err != nil {
			log.Debug().Msgf("error closing mongo cursor %v", err)
		}
	}(extraDocumentCursor, ctx)

	cursor, err := p.col(models.ProductCollectionName).Find(updatedCtx, filter, pagingOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("error finding products: %w", err)
	}
	products = make([]models.Product, cursor.RemainingBatchLength())
	if err = cursor.All(ctx, &products); err != nil {
		return nil, 0, fmt.Errorf("error getting remaining bacth products: %w", err)
	}

	return products, uint64(totalRecord), nil
}
