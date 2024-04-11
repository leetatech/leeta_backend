package interfaces

import (
	"github.com/leetatech/leeta_backend/pkg/query/filter"
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
