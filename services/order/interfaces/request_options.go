package interfaces

import (
	"github.com/leetatech/leeta_backend/pkg/query/filter"
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
