package interfaces

import (
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/sorting"
	"github.com/leetatech/leeta_backend/services/models"
)

var lgaRequestName = filter.ReadableValue[string]{
	Label: "LGA",
	Value: "lga",
}

var productIDRequestName = filter.ReadableValue[string]{
	Label: "Product",
	Value: "product_id",
}

var feeTypeRequestName = filter.ReadableValue[string]{
	Label: "Fee Type",
	Value: "fee_type",
}

var statusRequestName = filter.ReadableValue[string]{
	Label: "Status",
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
			Type: filter.ControlTypeString,
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
			Type: filter.ControlTypeEnum,
		},
		Operators: []filter.ReadableValue[filter.CompareOperator]{
			operatorEqual,
		},
		Values: []string{
			string(models.DeliveryFee),
			string(models.ServiceFee),
			string(models.ProductFee),
		},
	},
	{
		Name: statusRequestName,
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
