package domain

type Vendor struct {
	ID       string `json:"id" bson:"id"`
	Status   string `json:"status" bson:"status"`
	StatusTs int64  `json:"status_ts" bson:"status_ts"`
	Ts       string `json:"ts" bson:"ts"`
} // @name Vendor

type Customer struct {
	ID       string `json:"id" bson:"id"`
	Status   string `json:"status" bson:"status"`
	StatusTs int64  `json:"status_ts" bson:"status_ts"`
	Ts       string `json:"ts" bson:"ts"`
} // @name Customer
