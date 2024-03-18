package interfaces

import (
	"errors"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/product/domain"
	"net/http"
	"strconv"
)

func checkFormFileAndAddProducts(r *http.Request) (*domain.ProductRequest, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.FormParseError, errors.New("failed to parse multipart form"))
	}

	var (
		vendorId            = r.FormValue("vendor_id")
		parentCategory      = r.FormValue("parent_category")
		subCategory         = r.FormValue("sub_category")
		name                = r.FormValue("name")
		weight              = r.FormValue("weight")
		description         = r.FormValue("description")
		originalPrice       = r.FormValue("original_price")
		vat                 = r.FormValue("vat")
		originalPriceAndVat = r.FormValue("original_price_and_vat")
		discount            = r.FormValue("discount") == "true"
		discountPrice       = r.FormValue("discount_price")
		status              = r.FormValue("status")
	)

	images, err := GetImages(r)
	if err != nil {
		return nil, err
	}

	setParentCategory, err := models.SetProductCategory(models.ProductCategory(parentCategory))
	if err != nil {
		return nil, err
	}
	setSubCategory, err := models.SetProductSubCategory(models.ProductSubCategory(subCategory))
	if err != nil {
		return nil, err
	}

	price, err := stringToFloat64(originalPrice)
	if err != nil {
		return nil, err
	}

	vatPrice, err := stringToFloat64(vat)
	if err != nil {
		return nil, err
	}

	priceAndVat, err := stringToFloat64(originalPriceAndVat)
	if err != nil {
		return nil, err
	}

	totalDiscount, err := stringToFloat64(discountPrice)
	if err != nil {
		return nil, err
	}

	setStatus, err := models.SetProductStatus(models.ProductStatus(status))
	if err != nil {
		return nil, err
	}

	return &domain.ProductRequest{
		VendorID:            vendorId,
		ParentCategory:      setParentCategory,
		SubCategory:         setSubCategory,
		Images:              images,
		Name:                name,
		Weight:              weight,
		Description:         description,
		OriginalPrice:       price,
		Vat:                 vatPrice,
		OriginalPriceAndVat: priceAndVat,
		Discount:            discount,
		DiscountPrice:       totalDiscount,
		Status:              setStatus,
	}, nil
}

func GetImages(r *http.Request) ([]string, error) {

	var images []string
	files := r.MultipartForm.File["images"]
	if files == nil {
		return nil, leetError.ErrorResponseBody(leetError.UnmarshalError, errors.New("images field cannot be empty"))
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, leetError.ErrorResponseBody(leetError.FormParseError, errors.New("failed to get image from the request"))
		}
		defer file.Close()

		imageFormat, err := pkg.CheckImageFormat(fileHeader)
		if err != nil {
			return nil, err
		}

		img, err := pkg.CheckImageSizeAndDimension(fileHeader, file, 500, 600)
		if err != nil {
			return nil, err
		}

		encodedImage, err := pkg.EncodeImageToBase64(img, imageFormat)
		if err != nil {
			return nil, leetError.ErrorResponseBody(leetError.EncryptionError, err)
		}

		images = append(images, encodedImage)

		return images, nil
	}

	return nil, nil
}

func stringToFloat64(strValue string) (float64, error) {
	value, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return 0, leetError.ErrorResponseBody(leetError.UnmarshalError, err)
	}

	return value, nil
}

func ToFilterOption(options filter.RequestOption, _ int) filter.RequestOption {
	return options
}
