package models

import "errors"

type Verification struct {
	ID              string              `json:"id" bson:"id"`
	Code            string              `json:"code" bson:"code"`
	Topic           string              `json:"topic" bson:"topic"`
	Type            MessageDeliveryType `json:"type" bson:"type"`
	Target          string              `json:"target" bson:"target"`
	ExpiresAt       int64               `json:"expires_at" bson:"expires_at"`
	Validated       bool                `json:"validated" bson:"validated"`
	StatusTimeStamp int64               `json:"status_ts" bson:"status_ts"`
	Timestamp       int64               `json:"ts" bson:"ts"`
} // @name Verification

type MessageDeliveryType string

const (
	SMS   MessageDeliveryType = "SMS"
	EMAIL MessageDeliveryType = "EMAIL"
	PUSH  MessageDeliveryType = "PUSH"
)

func IsValidMessageDeliveryType(deliveryType MessageDeliveryType) bool {
	return deliveryType == SMS || deliveryType == EMAIL || deliveryType == PUSH
}

func SetIsValidMessageDeliveryType(deliveryType MessageDeliveryType) (MessageDeliveryType, error) {
	switch IsValidMessageDeliveryType(deliveryType) {
	case true:
		return deliveryType, nil
	default:
		return "", errors.New("invalid onboarding status")
	}
}
