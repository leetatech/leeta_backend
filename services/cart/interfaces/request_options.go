package interfaces

import (
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/sorting"
)

var cartIDRequestName = filter.ReadableValue[string]{
	Label: "cart id",
	Value: "id",
}

// LabelIsEqualTo holds filter request options operator labels
const (
	LabelIsEqualTo = "is equal to"
)

var operatorEqual = filter.ReadableValue[filter.CompareOperator]{
	Label: LabelIsEqualTo,
	Value: filter.CompareOperatorIsEqualTo,
}

var listCartOptions = []filter.RequestOption{
	{
		Name: cartIDRequestName,
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
