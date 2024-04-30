package interfaces

import (
	"github.com/leetatech/leeta_backend/pkg/query/filter"
)

var lgaRequestName = filter.ReadableValue[string]{
	Label: "lga",
	Value: "lga",
}

var productIDRequestName = filter.ReadableValue[string]{
	Label: "product_id",
	Value: "product_id",
}

var feeTypeRequestName = filter.ReadableValue[string]{
	Label: "fee_type",
	Value: "fee_type",
}

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

var listFeesOptions = []filter.RequestOption{
	{
		Name: lgaRequestName,
		Control: filter.RequestOptionType{
			Type: "LGA",
		},
		Operators: []filter.ReadableValue[filter.CompareOperator]{
			operatorEqual,
		},
		MultiSelect: true,
	},
	{
		Name: productIDRequestName,
		Control: filter.RequestOptionType{
			Type: filter.ControlTypeString,
		},
		Operators: []filter.ReadableValue[filter.CompareOperator]{
			operatorEqual,
		},
		MultiSelect: true,
	},
	{
		Name: feeTypeRequestName,
		Control: filter.RequestOptionType{
			Type: "FeeType",
		},
		Operators: []filter.ReadableValue[filter.CompareOperator]{
			operatorEqual,
		},
		MultiSelect: true,
	},
	{
		Name: statusRequestName,
		Control: filter.RequestOptionType{
			Type: "FeesStatuses",
		},
		Operators: []filter.ReadableValue[filter.CompareOperator]{
			operatorEqual,
		},
		MultiSelect: true,
	},
}
