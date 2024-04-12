package database

import (
	"context"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg/config"
	"github.com/leetatech/leeta_backend/pkg/query/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoDBClient(ctx context.Context, config *config.ServerConfig) (*mongo.Client, error) {
	clientOpts := config.GetClientOptions()
	mongoClient, err := mongo.NewClient(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("error initializing new mongo client %w", err)
	}
	err = mongoClient.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongo client: %w", err)
	}
	return mongoClient, nil
}

func GetPaginatedOpts(limit, page int64) *options.FindOptions {
	l := limit
	skip := page*limit - limit
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}

	return &fOpt
}

func BuildMongoFilterQuery(filter *filter.Request) bson.M {
	query := bson.M{}

	if filter == nil {
		return query
	}

	switch filter.Operator {
	case "and":
		for _, field := range filter.Fields {
			query[field.Name] = field.Value
		}
	case "or":
		var orConditions []bson.M
		for _, field := range filter.Fields {
			orConditions = append(orConditions, bson.M{field.Name: field.Value})
		}
		query["$or"] = orConditions
	}

	return query
}
