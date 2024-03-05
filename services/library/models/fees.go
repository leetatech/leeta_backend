package models

import (
	"errors"
	"github.com/leetatech/leeta_backend/services/library/leetError"
)

type Fees struct {
	ID         string       `json:"id" bson:"id"`
	ProductID  string       `json:"product_id" bson:"product_id"`
	CostPerKg  float64      `json:"cost_per_kg,omitempty" bson:"cost_per_kg"`
	CostPerQty float64      `json:"cost_per_qty,omitempty" bson:"cost_per_qty"`
	ServiceFee float64      `json:"service_fee" bson:"service_fee"`
	Status     CartStatuses `json:"status" bson:"status"`
	StatusTs   int64        `json:"status_ts" bson:"status_ts"`
	Ts         int64        `json:"ts" bson:"ts"`
} // @name Fees

type FeesStatuses string

const (
	FeesActive   FeesStatuses = "ACTIVE"   // fees has been created and active
	FeesInactive FeesStatuses = "INACTIVE" // fees has been inactivated
)

func IsValidFeesStatus(status FeesStatuses) bool {
	return status == FeesActive || status == FeesInactive
}

func SetFeesStatus(status FeesStatuses) (FeesStatuses, error) {
	switch IsValidFeesStatus(status) {
	case true:
		return status, nil
	default:
		return "", leetError.ErrorResponseBody(leetError.FeesStatusesError, errors.New("invalid fees status"))
	}
}
