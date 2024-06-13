package database

import (
	"context"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg/config"
	"github.com/leetatech/leeta_backend/pkg/query/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
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
			query = andLogicOperator(query, field)
		}
	case "or":
		var orConditions []bson.M
		for _, field := range requestFilter.Fields {
			orConditions = orLogicOperator(field)
		}
		query["$or"] = orConditions
	}
	return query
}

func andLogicOperator(query bson.M, field filter.RequestField) bson.M {
	if reflect.TypeOf(field.Value).Kind() == reflect.Slice {
		switch field.Operator {
		case filter.CompareOperatorContains:
			query[field.Name] = bson.M{"$in": field.Value}
		case filter.CompareOperatorIsEqualTo:
			if values, isSlice := field.Value.([]any); isSlice {
				query[field.Name] = bson.M{"$eq": field.Value}
			} else {
				query[field.Name] = values[0]
			}
		}
	} else {
		query[field.Name] = field.Value
	}

	return query
}

func orLogicOperator(field filter.RequestField) []bson.M {
	var orConditions []bson.M
	if reflect.TypeOf(field.Value).Kind() == reflect.Slice {
		switch field.Operator {
		case filter.CompareOperatorContains:
			orConditions = append(orConditions, bson.M{field.Name: bson.M{"$in": field.Value}})
		case filter.CompareOperatorIsEqualTo:
			if values, isSlice := field.Value.([]any); isSlice {
				orConditions = append(orConditions, bson.M{field.Name: bson.M{"$eq": field.Value}})
			} else {
				orConditions = append(orConditions, bson.M{field.Name: values[0]})
			}
			orConditions = append(orConditions, bson.M{field.Name: field.Value.([]interface{})[0]})
		}
	} else {
		orConditions = append(orConditions, bson.M{field.Name: field.Value})
	}

	return orConditions
}
