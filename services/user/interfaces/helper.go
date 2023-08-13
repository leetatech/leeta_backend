package interfaces

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/models"
	"github.com/leetatech/leeta_backend/services/user/domain"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"
	"strings"
)

func encodeImageToBase64(img image.Image, format string) (string, error) {
	var buf bytes.Buffer

	switch format {
	case "jpeg":
		err := jpeg.Encode(&buf, img, nil)
		if err != nil {
			return "", err
		}
	case "png":
		err := png.Encode(&buf, img)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("unsupported image format")
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func checkFormFileSpecification(r *http.Request) (*domain.VendorVerificationRequest, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.FormParseError, errors.New("failed to parse multipart form"))
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
		return nil, leetError.ErrorResponseBody(leetError.FormParseError, errors.New("failed to get image from the request"))
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	imageFormat := ""
	switch {
	case strings.HasPrefix(contentType, "image/jpeg"):
		imageFormat = "jpeg"

	case strings.HasPrefix(contentType, "image/png"):
		imageFormat = "png"

	default:
		return nil, leetError.ErrorResponseBody(leetError.FormParseError, errors.New("invalid image format. Only JPEG and PNG images are allowed"))
	}

	const maxImageSize = 5 * 1024 * 1024 // 5MB
	if header.Size > maxImageSize {
		return nil, leetError.ErrorResponseBody(leetError.FormParseError, errors.New("image size exceeds the maximum limit of 5MB"))
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.FormParseError, errors.New("failed to decode the image"))
	}

	minWidth := 800
	minHeight := 800
	if img.Bounds().Dx() < minWidth || img.Bounds().Dy() < minHeight {
		return nil, leetError.ErrorResponseBody(leetError.FormParseError, errors.New("image dimensions should be at least 800x800 pixels"))
	}

	encodedImage, err := encodeImageToBase64(img, imageFormat)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.EncryptionError, err)
	}

	businessCategory, err := models.SetBusinessCategory(models.BusinessCategory(category))
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.FormParseError, err)
	}
	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.UnmarshalError, err)
	}
	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.UnmarshalError, err)
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
