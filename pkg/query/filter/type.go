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
)

const (
	CompareOperatorIsEqualTo CompareOperator = "isEqualTo"
	// CompareOperatorContains is a CompareOperator of type contains.
	CompareOperatorContains CompareOperator = "contains"
)

type ControlType string // @name ControlType

/*
CompareOperator ENUM(

	isEqualTo

)
*/
type CompareOperator string // @name CompareOperator
