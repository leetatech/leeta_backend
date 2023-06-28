package models

import (
	"errors"
)

type Identity struct {
	ID          string        `json:"id" bson:"id"`
	CustomerID  string        `json:"customer_id" bson:"customer_id"`
	Role        UserCategory  `json:"role" bson:"role"`
	Credentials []Credentials `json:"credentials" bson:"credentials"`
} // @name Identity

type Credentials struct {
	Type            CredentialType   `json:"type" bson:"type"`
	Password        string           `json:"password" bson:"password"`
	Status          CredentialStatus `json:"status" bson:"status"`
	StatusTimestamp int64            `json:"status_ts" bson:"status_ts"`
	Timestamp       int64            `json:"ts" bson:"ts"`
} // @name Credentials

/*
**constants/enums
 */

type CredentialType string

const (
	CredentialsTypeLogin CredentialType = "LOGIN"
	CredentialsTypePIN   CredentialType = "PIN"
)

type CredentialStatus string

const (
	CredentialStatusActive   CredentialStatus = "ACTIVE"
	CredentialStatusInactive CredentialStatus = "INACTIVE"
	CredentialStatusLocked   CredentialStatus = "LOCKED"
)

type UserCategory string

const (
	VendorCategory UserCategory = "vendor"
	BuyerCategory  UserCategory = "buyer"
)

func IsValidCredentialType(credentialType CredentialType) bool {
	return credentialType == CredentialsTypeLogin || credentialType == CredentialsTypePIN
}
func SetCredentialType(credentialType CredentialType) (CredentialType, error) {
	switch IsValidCredentialType(credentialType) {
	case true:
		return credentialType, nil
	default:
		return "", errors.New("invalid credential type")
	}
}

func IsValidCredentialStatus(status CredentialStatus) bool {
	return status == CredentialStatusActive || status == CredentialStatusInactive || status == CredentialStatusLocked
}
func SetCredentialStatus(status CredentialStatus) (CredentialStatus, error) {
	switch IsValidCredentialStatus(status) {
	case true:
		return status, nil
	default:
		return "", errors.New("invalid credential status")
	}
}

func IsValidUserCategory(category UserCategory) bool {
	return category == VendorCategory || category == BuyerCategory
}
func SetUserCategory(category UserCategory) (UserCategory, error) {
	switch IsValidUserCategory(category) {
	case true:
		return category, nil
	default:
		return "", errors.New("invalid user category")
	}
}
