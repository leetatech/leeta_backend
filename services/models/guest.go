package models

type Guest struct {
	ID                     string      `json:"id" bson:"id"`
	Location               Coordinates `json:"location,omitempty" bson:"location"`
	DeviceID               string      `json:"device_id,omitempty" bson:"device_id"`
	FirstName              string      `json:"first_name,omitempty" bson:"first_name"`
	LastName               string      `json:"last_name,omitempty" bson:"last_name"`
	Number                 string      `json:"number,omitempty" bson:"number"`
	Email                  string      `json:"email,omitempty" bson:"email"`
	Address                Address     `json:"address,omitempty" bson:"address"`
	DefaultDeliveryAddress bool        `json:"default_delivery_address,omitempty" bson:"default_delivery_address"`
}
