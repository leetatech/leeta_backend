package filter

// FilterRequest is a struct representing a filter request.
// Operator is the logic operator used for the request.
// Fields is a slice of RequestField, representing the fields to be used for the filtering.
type FilterRequest struct {
	Operator LogicOperator  `json:"operator" binding:"required"`
	Fields   []RequestField `json:"fields" binding:"dive"`
} // @name FilterRequest

type PagingRequest struct {
	PageIndex int `json:"index"`
	PageSize  int `json:"size"`
} // @name PagingRequest

/*
LogicOperator ENUM(

	and
	or

)
*/
type LogicOperator string // @name LogicOperator

/*
ControlType ENUM(

	enum
	float
	integer
	string
	dateTime
	uuid
	autocomplete

)
*/

const (
	// ControlTypeString is a ControlType of type string.
	ControlTypeString ControlType = "string"
)

const (
	CompareOperatorIsEqualTo CompareOperator = "isEqualTo"
)

type ControlType string // @name ControlType

/*
CompareOperator ENUM(

	isEqualTo

)
*/
type CompareOperator string // @name CompareOperator

// RequestField represents a field in a request
// Field Name: The name of the field
// Field Value: The value of the field, which can be a list of values or a single value
type RequestField struct {
	Name string `json:"name" binding:"required"`
	// Value can be a list of values or a value
	Value any `json:"value" binding:"required"`
} // @name RequestField

// ResultSelector is a type that represents the selection criteria for querying data. It contains a filter, sorting, and paging information.
// Filter is a pointer to a filter.Request struct that specifies the filtering criteria for the query.
// Sorting is a pointer to a sorting.Request struct that specifies the sorting order for the query.
// Paging is a pointer to a paging.Request struct that specifies the paging configuration for the query.
type ResultSelector struct {
	Filter *FilterRequest `json:"filter" binding:"omitempty"`
	Paging *PagingRequest `json:"paging" binding:"omitempty"`
} // @name ResultSelector

// ReadableValue is a generic type that represents a human-readable value with a corresponding backend value.
// It has two fields: `Label` (the human-readable form of the value) and `Value` (the value for the backend).
type ReadableValue[T any] struct {
	// Label is the human-readable form of the value
	Label string `json:"label"`
	// Value is the value for the backend
	Value T `json:"value"`
} // @name ReadableValue

// RequestOptionType configures the type of control for a field in a request option.
type RequestOptionType struct {
	Type ControlType `json:"type" enums:"string,float,integer,enum"`
} // @name RequestOptionType

// RequestOption configures a field for validation
//
// Name: The name of the option
// Control: The type of control for the option
// Operators: The list of comparison operators for the option
// Values: The possible values for the option
// MultiSelect: Indicates whether the option supports multiple selections
type RequestOption struct {
	Name        ReadableValue[string]
	Control     RequestOptionType
	Operators   []ReadableValue[CompareOperator]
	Values      []string
	MultiSelect bool
} // @name RequestOption
