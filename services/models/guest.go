package models

type Guest struct {
	ID       string      `json:"id" bson:"id"`
	Location Coordinates `json:"location" bson:"location"`
	DeviceID string      `json:"device_id" bson:"device_id"`
}
