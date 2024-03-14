package interfaces

import "github.com/leetatech/leeta_backend/services/library"

var productStatusRequestName = library.ReadableValue[string]{
	Label: "Product Status",
	Value: "product_status",
}

// Filter request options operator labels
const (
	LabelIsEqualTo = "is equal to"
)

var operatorEqual = library.ReadableValue[library.CompareOperator]{
	Label: LabelIsEqualTo,
	Value: library.CompareOperatorIsEqualTo,
}

var listProductOptions = []library.RequestOption{
	{
		Name: productStatusRequestName,
		Control: library.RequestOptionType{
			Type: library.ControlTypeString,
		},
		Operators: []library.ReadableValue[library.CompareOperator]{
			operatorEqual,
		},
		MultiSelect: true,
	},
}
