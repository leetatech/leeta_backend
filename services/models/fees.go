package models

type Fee struct {
	ID         string       `json:"id" bson:"id"`
	ProductID  string       `json:"product_id" bson:"product_id"`
	CostPerKg  float64      `json:"cost_per_kg,omitempty" bson:"cost_per_kg"`
	CostPerQty float64      `json:"cost_per_qty,omitempty" bson:"cost_per_qty"`
	ServiceFee float64      `json:"service_fee" bson:"service_fee"`
	Status     CartStatuses `json:"status" bson:"status"`
	StatusTs   int64        `json:"status_ts" bson:"status_ts"`
	Ts         int64        `json:"ts" bson:"ts"`
} // @name Fee

type FeesStatuses string

const (
	FeesActive   FeesStatuses = "ACTIVE"   // fees has been created and active
	FeesInactive FeesStatuses = "INACTIVE" // fees has been inactivated
)
