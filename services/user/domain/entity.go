package domain

import "github.com/leetatech/leeta_backend/services/library/models"

type VendorDetailsUpdateRequest struct {
	ID        string          `json:"id" bson:"_id"`
	Identity  string          `json:"identity" bson:"identity"`
	FirstName string          `json:"first_name" bson:"first_name"`
	LastName  string          `json:"last_name" bson:"last_name"`
	Status    models.Statuses `json:"status" bson:"status"`
} // @name VendorDetailsUpdateRequest

type VendorVerificationRequest struct {
	FirstName   string                  `json:"first_name" bson:"first_name"`
	LastName    string                  `json:"last_name" bson:"last_name"`
	Identity    string                  `json:"identity" bson:"identity"`
	Name        string                  `json:"name" bson:"name"`
	CAC         string                  `json:"cac" bson:"cac"`
	Category    models.BusinessCategory `json:"category" bson:"category"`
	Description string                  `json:"description" bson:"description"`
	Phone       []models.Phone          `json:"phone" bson:"phone"`
	Address     []models.Address        `json:"address" bson:"address"`
} // @name VendorVerificationRequest
