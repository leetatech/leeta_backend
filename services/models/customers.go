package models

import (
	"errors"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"strings"
)

type User struct {
	ID            string   `json:"id" bson:"id"`
	FirstName     string   `json:"first_name" bson:"first_name"`
	LastName      string   `json:"last_name" bson:"last_name"`
	Email         Email    `json:"email" bson:"email"`
	Address       Address  `json:"address" bson:"address"`
	Phone         Phone    `json:"phone" bson:"phone"`
	DOB           string   `json:"dob" bson:"dob"`
	HasPIN        bool     `json:"has_pin" bson:"has_pin"`
	PinBlocked    bool     `json:"pin_blocked" bson:"pin_blocked"`
	IsBlocked     bool     `json:"is_blocked" bson:"is_blocked"`
	BlockedReason string   `json:"is_blocked_reason" bson:"is_blocked_reason"`
	Status        Statuses `json:"status" bson:"status"`
}

func (user *User) ExtractName(fullName string) error {
	names := strings.Fields(fullName)

	// Handle cases where there is only one name or more than two names
	switch len(names) {
	case 0:
		return errors.New("no user names provides")
	case 1:
		// Only one name provided, consider it as the first name
		user.FirstName = names[0]
		return nil
	default:
		// More than one name provided, consider the first as the first name
		// and the rest as the last name
		user.FirstName = names[0]
		user.LastName = strings.Join(names[1:], " ")
		return nil
	}
}

type Customer struct {
	User
	TimeStamps
} // @name Customer

type TimeStamps struct {
	StatusTime int64 `json:"status_ts" bson:"status_ts"`
	Time       int64 `json:"ts" bson:"ts"`
}

type Vendor struct {
	User
	AdminID string `json:"admin_id" bson:"admin_id"`
	TimeStamps
} // @name Vendor

type Admin struct {
	User
	Department string `json:"department"`
	Role       string `json:"role"`
	TimeStamps
} // @name Admin

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
	State           string      `json:"state" bson:"state"`
	City            string      `json:"city" bson:"city"`
	LGA             string      `json:"lga" bson:"lga"`
	FullAddress     string      `json:"full_address" bson:"full_address"`
	ClosestLandmark string      `json:"closest_landmark" bson:"closest_landmark"`
	Coordinates     Coordinates `json:"coordinate" bson:"coordinate"`
	Verified        bool        `json:"verified" bson:"verified"`
} // @name Address

// Coordinates model
type Coordinates struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
} // @name Coordinates

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
