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

// BuildMongoFilterQuery constructs a MongoDB filter query based on the provided request filter
// and a field mapping. It supports both "and" and "or" operators for combining field conditions.
// TODO: handle filter field mapping better
func BuildMongoFilterQuery(requestFilter *filter.Request, fieldMapping map[string]string) bson.M {
	query := bson.M{}

	if requestFilter == nil {
		return query
	}

	// Helper function to build individual field queries
	buildFieldQuery := func(field filter.RequestField) bson.M {
		// Use the mapped field name if it exists, otherwise use the original field name
		fieldName := fieldMapping[field.Name]
		if fieldName == "" {
			fieldName = field.Name
		}

		if field.Operator == filter.CompareOperatorContains {
			return bson.M{fieldName: bson.M{"$in": field.Value}}
		}
		return bson.M{fieldName: field.Value}
	}

	switch requestFilter.Operator {
	case "and":
		for _, field := range requestFilter.Fields {
			fieldQuery := buildFieldQuery(field)
			for key, value := range fieldQuery {
				query[key] = value
			}
		}
	case "or":
		orConditions := make([]bson.M, len(requestFilter.Fields))
		for i, field := range requestFilter.Fields {
			orConditions[i] = buildFieldQuery(field)
		}
		query["$or"] = orConditions
	}

	return query
}
