package interfaces

import (
	"errors"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/pkg/helpers"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/user/domain"
	"net/http"
	"strconv"
)

func checkFormFileSpecification(r *http.Request) (*domain.VendorVerificationRequest, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, errs.Body(errs.FormParseError, errors.New("failed to parse multipart form"))
	}

	var (
		firstName       = r.FormValue("first_name")
		lastName        = r.FormValue("last_name")
		businessName    = r.FormValue("business_name")
		cac             = r.FormValue("cac")
		category        = r.FormValue("business_category")
		description     = r.FormValue("description")
		primaryPhone    = r.FormValue("primary_phone") == "true"
		phoneNumber     = r.FormValue("phone_number")
		State           = r.FormValue("state")
		city            = r.FormValue("city")
		lga             = r.FormValue("lga")
		fullAddress     = r.FormValue("full_address")
		closestLandmark = r.FormValue("closest_landmark")
		latitudeStr     = r.FormValue("latitude")
		longitudeStr    = r.FormValue("longitude")
	)

	file, header, err := r.FormFile("image")
	if err != nil {
		return nil, errs.Body(errs.FormParseError, errors.New("failed to get image from the request"))
	}
	defer file.Close()

	imageFormat, err := helpers.CheckImageFormat(header)
	if err != nil {
		return nil, err
	}

	img, err := helpers.CheckImageSizeAndDimension(header, file, 800, 800)
	if err != nil {
		return nil, err
	}

	encodedImage, err := helpers.EncodeImageToBase64(img, imageFormat)
	if err != nil {
		return nil, errs.Body(errs.EncryptionError, err)
	}

	businessCategory, err := models.SetBusinessCategory(models.BusinessCategory(category))
	if err != nil {
		return nil, errs.Body(errs.FormParseError, err)
	}
	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		return nil, errs.Body(errs.UnmarshalError, err)
	}
	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		return nil, errs.Body(errs.UnmarshalError, err)
	}

	request := domain.VendorVerificationRequest{
		FirstName:   firstName,
		LastName:    lastName,
		Identity:    encodedImage,
		Name:        businessName,
		CAC:         cac,
		Category:    businessCategory,
		Description: description,
		Phone: []models.Phone{
			{
				Primary: primaryPhone,
				Number:  phoneNumber,
			},
		},
		Address: []models.Address{
			{
				State:           State,
				City:            city,
				LGA:             lga,
				FullAddress:     fullAddress,
				ClosestLandmark: closestLandmark,
				Coordinates: models.Coordinates{
					Latitude:  latitude,
					Longitude: longitude,
				},
			},
		},
	}

	return &request, nil
}
