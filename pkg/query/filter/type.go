package filter

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
	// CompareOperatorContains is a CompareOperator of type contains.
	CompareOperatorContains CompareOperator = "contains"
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
