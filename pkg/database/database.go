package database

import (
	"context"
	"fmt"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/leetatech/leeta_backend/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoDBClient(ctx context.Context, config *config.ServerConfig) (*mongo.Client, error) {
	clientOpts := config.GetClientOptions()
	mongoClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("error initializing new mongo client %w", err)
	}
	return mongoClient, nil
}

func GetPaginatedOpts(limit, page int64) *options.FindOptions {
	l := limit
	skip := page*limit - limit
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}

	return &fOpt
}

func BuildMongoFilterQuery(requestFilter *filter.Request) bson.M {
	query := bson.M{}

	if requestFilter == nil {
		return query
	}

	switch requestFilter.Operator {
	case "and":
		for _, field := range requestFilter.Fields {
			if field.Operator == filter.CompareOperatorContains {
				query[field.Name] = bson.M{"$in": field.Value}
			} else {
				query[field.Name] = field.Value
			}
		}
	case "or":
		var orConditions []bson.M
		for _, field := range requestFilter.Fields {
			if field.Operator == filter.CompareOperatorContains {
				orConditions = append(orConditions, bson.M{field.Name: bson.M{"$in": field.Value}})
			} else {
				orConditions = append(orConditions, bson.M{field.Name: field.Value})
			}
		}
		query["$or"] = orConditions
	}

	return query
}
