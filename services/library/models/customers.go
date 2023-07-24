package models

import (
	"errors"
	"github.com/leetatech/leeta_backend/services/library/leetError"
)

type Vendor struct {
	ID              string   `json:"id" bson:"id"`
	FirstName       string   `json:"first_name" bson:"first_name"`
	LastName        string   `json:"last_name" bson:"last_name"`
	Dob             string   `json:"dob" bson:"dob"`
	Email           Email    `json:"email" bson:"email"`
	HasPIN          bool     `json:"has_pin" bson:"has_pin"`
	PinBlocked      bool     `json:"pin_blocked" bson:"pin_blocked"`
	IsBlocked       bool     `json:"is_blocked" bson:"is_blocked"`
	BlockedReason   string   `json:"is_blocked_reason" bson:"is_blocked_reason"`
	Status          Statuses `json:"status" bson:"status"`
	StatusTimeStamp int64    `json:"status_ts" bson:"status_ts"`
	Timestamp       int64    `json:"ts" bson:"ts"`
} // @name Vendor

// Phone model
type Phone struct {
	Primary  bool   `json:"primary" bson:"primary"`
	Number   string `json:"number" bson:"number"`
	Verified bool   `json:"verified" bson:"verified"`
} // @name Phone

// Email model
type Email struct {
	Address  string `json:"address" bson:"address"`
	Verified bool   `json:"verified" bson:"verified"`
} // @name Email

// Address model
type Address struct {
	State           string `json:"state" bson:"state"`
	City            string `json:"city" bson:"city"`
	LGA             string `json:"lga" bson:"lga"`
	FullAddress     string `json:"full_address" bson:"full_address"`
	ClosestLandmark string `json:"closest_landmark" bson:"closest_landmark"`
	Verified        bool   `json:"verified" bson:"verified"`
} // @name Address

// Business - vendor business
type Business struct {
	ID              string           `json:"id" bson:"id"`
	VendorID        string           `json:"vendor_id" bson:"vendor_id"`
	Name            string           `json:"name" bson:"name"`
	CAC             string           `json:"cac" bson:"cac"`
	Category        BusinessCategory `json:"category" bson:"category"`
	Description     string           `json:"description" bson:"description"`
	Phone           []Phone          `json:"phone" bson:"phone"`
	Address         []Address        `json:"address" bson:"address"`
	Status          Statuses         `json:"status" bson:"status"`
	StatusTimeStamp int64            `json:"status_ts" bson:"status_ts"`
	Timestamp       int64            `json:"ts" bson:"ts"`
} // @name Business

/*
**constants/enums
 */

// Statuses type
type Statuses string

const (
	SignedUp   Statuses = "SIGNEDUP"   // just signed up
	Registered Statuses = "REGISTERED" // filled the required information
	Verified   Statuses = "VERIFIED"   // all details verified
	Onboarded  Statuses = "ONBOARDED"  // now fully onboarded
	Rejected   Statuses = "REJECTED"   // rejected
	Exited     Statuses = "EXITED"     // no longer exists
	Locked     Statuses = "LOCKED"     // currently locked for some reasons
)

// BusinessCategory type
type BusinessCategory string

const (
	LPG BusinessCategory = "LPG"
	LNG BusinessCategory = "LPG"
)

func IsValidStatuses(status Statuses) bool {
	return status == SignedUp || status == Registered || status == Verified || status == Onboarded || status == Rejected || status == Exited || status == Locked
}

func SetIsValidStatuses(status Statuses) (Statuses, error) {
	switch IsValidStatuses(status) {
	case true:
		return status, nil
	default:
		return "", leetError.ErrorResponseBody(leetError.BusinessCategoryError, errors.New("invalid business category"))
	}
}

func IsValidBusinessCategory(category BusinessCategory) bool {
	return category == LPG || category == LNG
}

func SetBusinessCategory(category BusinessCategory) (BusinessCategory, error) {
	switch IsValidBusinessCategory(category) {
	case true:
		return category, nil
	default:
		return "", leetError.ErrorResponseBody(leetError.BusinessCategoryError, errors.New("invalid business category"))
	}
}
