package domain

import "github.com/leetatech/leeta_backend/services/library/models"

type GetVendorProductsRequest struct {
	VendorID      string                 `json:"vendor_id" bson:"vendor_id"`
	ProductStatus []models.ProductStatus `json:"product_status" bson:"product_status"`
	Limit         int64                  `json:"limit" bson:"limit"`
	Page          int64                  `json:"page" bson:"page"`
} // @name GetVendorProductRequest

type GetVendorProductsResponse struct {
	Products    []models.Product `json:"products" bson:"products"`
	HasNextPage bool             `json:"has_next_page" bson:"has_next_page"`
}

type ProductRequest struct {
	VendorID            string                    `json:"vendor_id"`
	ParentCategory      models.ProductCategory    `json:"parent_category"`
	SubCategory         models.ProductSubCategory `json:"sub_category"`
	Images              []string                  `json:"images"`
	Name                string                    `json:"name"`
	Weight              string                    `json:"weight"`
	Description         string                    `json:"description"`
	OriginalPrice       float64                   `json:"original_price"`
	Vat                 float64                   `json:"vat"`
	OriginalPriceAndVat float64                   `json:"original_price_and_vat"`
	Discount            bool                      `json:"discount"`
	DiscountPrice       float64                   `json:"discount_price"`
	Status              models.ProductStatus      `json:"status"`
}

type GasProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
