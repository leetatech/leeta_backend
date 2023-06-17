package domain

import (
	"github.com/gofrs/uuid"
)

type Order struct {
	ID         uuid.UUID `json:"id" sql:"id"`
	CustomerID uuid.UUID `json:"source_user_id" sql:"source_user_id"`
	VendorID   uuid.UUID `json:"source_account_id" sql:"source_account_id"`
	Status     string    `json:"status" bson:"status"`
	StatusTs   int64     `json:"status_ts" bson:"status_ts"`
	Ts         string    `json:"ts" bson:"ts"`
} // @name Order
