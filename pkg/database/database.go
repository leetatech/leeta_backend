package database

import (
	"context"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg/config"
	"github.com/leetatech/leeta_backend/pkg/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func MongoDBClient(ctx context.Context, config *config.ServerConfig) (*mongo.Client, error) {
	clientOpts := config.GetClientOptions()
	mongoClient, err := mongo.NewClient(clientOpts)
	err = mongoClient.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongo client: %w", err)
	}
	return mongoClient, nil
}

func BuildMongoFilterQuery(filter *filter.FilterRequest) bson.M {
	query := bson.M{}

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
