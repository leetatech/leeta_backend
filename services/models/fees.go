package models

type Fee struct {
	ID        string       `json:"id" bson:"id"`
	ProductID string       `json:"product_id" bson:"product_id"`
	FeeType   FeeType      `json:"fee_type,omitempty" bson:"fee_type"`
	LGA       LGA          `json:"lga" bson:"lga"`
	Cost      Cost         `json:"cost" bson:"cost"`
	Status    FeesStatuses `json:"status" bson:"status"`
	StatusTs  int64        `json:"status_ts" bson:"status_ts"`
	Ts        int64        `json:"ts" bson:"ts"`
} // @name Fee

type FeesStatuses string

const (
	FeesActive   FeesStatuses = "ACTIVE"   // fees has been created and active
	FeesInactive FeesStatuses = "INACTIVE" // fees has been inactivated
)

type FeeType string

const (
	ServiceFee  FeeType = "SERVICE_FEE"
	ProductFee  FeeType = "PRODUCT_FEE"
	DeliveryFee FeeType = "DELIVERY_FEE"
)

type Cost struct {
	CostPerKG   float64 `json:"cost_per_kg" bson:"cost_per_kg"`
	CostPerQt   float64 `json:"cost_per_qty" bson:"cost_per_qty"`
	CostPerType float64 `json:"cost_per_type" bson:"cost_per_type"`
}

type LGA struct {
	State string `json:"state" bson:"state"`
	LGA   string `json:"lga" bson:"lga"`
}
