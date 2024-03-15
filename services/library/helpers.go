package library

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/leetatech/leeta_backend/services/library/filter"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"strings"
)

func CheckErrorType(err error, w http.ResponseWriter) {
	switch err := err.(type) {
	case *leetError.ErrorResponse:
		if err.Code() == leetError.ErrorUnauthorized {
			EncodeErrorResult(w, http.StatusUnauthorized)
			return
		}
	default:
		EncodeErrorResult(w, http.StatusInternalServerError)
		return
	}

}

func GetPaginatedOpts(limit, page int64) *options.FindOptions {
	l := limit
	skip := page*limit - limit
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}

	return &fOpt
}

func EncodeImageToBase64(img image.Image, format string) (string, error) {
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

func CheckImageFormat(fileHeader *multipart.FileHeader) (string, error) {
	contentType := fileHeader.Header.Get("Content-Type")
	switch {
	case strings.HasPrefix(contentType, "image/jpeg"):
		return "jpeg", nil

	case strings.HasPrefix(contentType, "image/png"):
		return "png", nil

	default:
		return "", leetError.ErrorResponseBody(leetError.FormParseError, errors.New("invalid image format. Only JPEG and PNG images are allowed"))
	}

}

func CheckImageSizeAndDimension(fileHeader *multipart.FileHeader, file multipart.File, width, height int) (image.Image, error) {
	const maxImageSize = 5 * 1024 * 1024 // 5MB
	if fileHeader.Size > maxImageSize {
		return nil, leetError.ErrorResponseBody(leetError.FormParseError, errors.New("image size exceeds the maximum limit of 5MB"))
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.FormParseError, errors.New("failed to decode the image"))
	}

	if img.Bounds().Dx() < width || img.Bounds().Dy() < height {
		return nil, leetError.ErrorResponseBody(leetError.FormParseError, errors.New("image dimensions should be at least 500x600 pixels"))
	}

	return img, nil
}

func BuildMongoFilterQuery(filter *filter.FilterRequest) bson.M {
	query := bson.M{}

	switch filter.Operator {
	case "and":
		for _, field := range filter.Fields {
			query[field.Name] = field.Value
		}
	case "or":
		var orConditions []bson.M
		for _, field := range filter.Fields {
			orConditions = append(orConditions, bson.M{field.Name: field.Value})
		}
		query["$or"] = orConditions
	}

	return query
}
