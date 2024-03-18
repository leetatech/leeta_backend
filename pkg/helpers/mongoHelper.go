package helpers

import (
	"github.com/leetatech/leeta_backend/pkg/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPaginatedOpts(limit, page int64) *options.FindOptions {
	l := limit
	skip := page*limit - limit
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}

	return &fOpt
}

func BuildMongoFilterQuery(filter *filter.FilterRequest) bson.M {
	query := bson.M{}

	if filter == nil {
		return bson.M{}
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
