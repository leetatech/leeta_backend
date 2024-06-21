package interfaces

import (
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/sorting"
)

var statusRequestName = filter.ReadableValue[string]{
	Label: "status",
	Value: "status",
}

// LabelIsEqualTo holds filter request options operator labels
const (
	LabelIsEqualTo = "is equal to"
)

var operatorEqual = filter.ReadableValue[filter.CompareOperator]{
	Label: LabelIsEqualTo,
	Value: filter.CompareOperatorIsEqualTo,
}

var listOrdersOptions = []filter.RequestOption{
	{
		Name: statusRequestName,
		Control: filter.RequestOptionType{
			Type: "OrderStatuses",
		},
		Operators: []filter.ReadableValue[filter.CompareOperator]{
			operatorEqual,
		},
		MultiSelect: true,
	},
}

var allowedSortFields = []string{"name"}

var defaultSortingRequest = &sorting.Request{
	SortColumn:    "name",
	SortDirection: sorting.DirectionDescending,
}
