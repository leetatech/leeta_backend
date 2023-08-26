package models

import (
	"errors"
	"github.com/leetatech/leeta_backend/services/library/leetError"
)

type Product struct {
	ID                  string           `json:"id" bson:"id"`
	ParentCategory      BusinessCategory `json:"parent_category" bson:"parent_category"`
	SubCategory         string           `json:"sub_category" bson:"sub_category"`
	Images              []string         `json:"images" bson:"images"`
	Name                string           `json:"name" bson:"name"`
	Weight              string           `json:"weight" bson:"weight"`
	Description         string           `json:"description" bson:"description"`
	OriginalPrice       float64          `json:"original_price" bson:"original_price"`
	Vat                 float64          `json:"vat" bson:"vat"`
	OriginalPriceAndVat float64          `json:"original_price_and_vat" bson:"original_price_and_vat"`
	Discount            bool             `json:"discount" bson:"discount"`
	DiscountPrice       float64          `json:"discount_price" bson:"discount_price"`
	FinalPrice          float64          `json:"final_price" bson:"final_price"`
	Status              string           `json:"status" bson:"status"`
	StatusTs            int64            `json:"status_ts" bson:"status_ts"`
	Ts                  int64            `json:"ts" bson:"ts"`
}

// ProductCategory type
type ProductCategory string

const (
	LPGProductCategory ProductCategory = "LPG"
	LNGProductCategory ProductCategory = "LPG"
)

// ProductSubCategory type
type ProductSubCategory string

const (
	CylinderSubCategory    ProductSubCategory = "GAS CYLINDER"
	CookerSubCategory      ProductSubCategory = "GAS COOKER"
	AccessoriesSubCategory ProductSubCategory = "GAS ACCESSORIES"
)

func IsValidProductCategory(category ProductCategory) bool {
	return category == LPGProductCategory || category == LNGProductCategory
}

func SetProductCategory(category ProductCategory) (ProductCategory, error) {
	switch IsValidProductCategory(category) {
	case true:
		return category, nil
	default:
		return "", leetError.ErrorResponseBody(leetError.BusinessCategoryError, errors.New("invalid business category"))
	}
}

func IsValidProductSubCategory(category ProductSubCategory) bool {
	return category == CylinderSubCategory || category == CookerSubCategory || category == AccessoriesSubCategory
}

func SetProductSubCategory(category ProductSubCategory) (ProductSubCategory, error) {
	switch IsValidProductSubCategory(category) {
	case true:
		return category, nil
	default:
		return "", leetError.ErrorResponseBody(leetError.BusinessCategoryError, errors.New("invalid business category"))
	}
}
