package models

import (
	"errors"
	"github.com/leetatech/leeta_backend/pkg/leetError"
)

type Product struct {
	ID                  string             `json:"id,omitempty" bson:"id"`
	VendorID            string             `json:"vendor_id,omitempty" bson:"vendor_id"`
	ParentCategory      ProductCategory    `json:"parent_category,omitempty" bson:"parent_category"`
	SubCategory         ProductSubCategory `json:"sub_category,omitempty" bson:"sub_category"`
	Images              []string           `json:"images,omitempty" bson:"images"`
	Name                string             `json:"name,omitempty" bson:"name"`
	Weight              string             `json:"weight,omitempty" bson:"weight"`
	Description         string             `json:"description,omitempty" bson:"description"`
	OriginalPrice       float64            `json:"original_price,omitempty" bson:"original_price"`
	Vat                 float64            `json:"vat,omitempty" bson:"vat"`
	OriginalPriceAndVat float64            `json:"original_price_and_vat,omitempty" bson:"original_price_and_vat"`
	Discount            bool               `json:"discount,omitempty" bson:"discount"`
	DiscountPrice       float64            `json:"discount_price,omitempty" bson:"discount_price"`
	FinalPrice          float64            `json:"final_price,omitempty" bson:"final_price"`
	Status              ProductStatus      `json:"status" bson:"status"`
	StatusTs            int64              `json:"status_ts" bson:"status_ts"`
	Ts                  int64              `json:"ts" bson:"ts"`
} // @name Product

// ProductCategory type
type ProductCategory string

const (
	LPGProductCategory ProductCategory = "LPG"
	LNGProductCategory ProductCategory = "LNG"
)

// ProductSubCategory type
type ProductSubCategory string

const (
	CylinderSubCategory    ProductSubCategory = "CYLINDER"
	CookerSubCategory      ProductSubCategory = "COOKER"
	AccessoriesSubCategory ProductSubCategory = "ACCESSORIES"
)

// ProductCategory type
type ProductStatus string

const (
	InStock    ProductStatus = "InStock"
	OutOfStock ProductStatus = "OutOfStock"
)

func IsValidProductCategory(category ProductCategory) bool {
	return category == LPGProductCategory || category == LNGProductCategory
}

func SetProductCategory(category ProductCategory) (ProductCategory, error) {
	switch IsValidProductCategory(category) {
	case true:
		return category, nil
	default:
		return "", leetError.ErrorResponseBody(leetError.ProductCategoryError, errors.New("invalid business category"))
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
		return "", leetError.ErrorResponseBody(leetError.ProductSubCategoryError, errors.New("invalid business category"))
	}
}

func IsValidProductStatus(status ProductStatus) bool {
	return status == InStock || status == OutOfStock
}

func SetProductStatus(status ProductStatus) (ProductStatus, error) {
	switch IsValidProductStatus(status) {
	case true:
		return status, nil
	default:
		return "", leetError.ErrorResponseBody(leetError.ProductStatusError, errors.New("invalid business category"))
	}
}
