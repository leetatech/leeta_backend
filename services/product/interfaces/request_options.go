package interfaces

import (
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/sorting"
)

var productStatusRequestName = filter.ReadableValue[string]{
	Label: "Product Status",
	Value: "status",
}

const (
	LabelIsEqualTo = "is equal to"
)

var operatorEqual = filter.ReadableValue[filter.CompareOperator]{
	Label: LabelIsEqualTo,
	Value: filter.CompareOperatorIsEqualTo,
}

var listProductOptions = []filter.RequestOption{
	{
		Name: productStatusRequestName,
		Control: filter.RequestOptionType{
			Type: filter.ControlTypeString,
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
