package models

type EarlyAccess struct {
	Email     string `json:"email" bson:"email"`
	Timestamp int64  `json:"ts" bson:"ts"`
} // @name EarlyAccess
