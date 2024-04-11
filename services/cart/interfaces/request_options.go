package interfaces

import (
	"github.com/leetatech/leeta_backend/pkg/query/filter"
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
