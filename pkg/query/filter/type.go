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
	models.LGA

)
*/

const (
	// ControlTypeString is a ControlType of type string.
	ControlTypeString ControlType = "string"
	ControlTypeLGA    ControlType = "models.LGA"
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
