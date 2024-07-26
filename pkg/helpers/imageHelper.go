package helpers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"strings"
)

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
		return "", errs.Body(errs.FormParseError, errors.New("invalid image format. Only JPEG and PNG images are allowed"))
	}

}

func CheckImageSizeAndDimension(fileHeader *multipart.FileHeader, file multipart.File, width, height int) (image.Image, error) {
	const maxImageSize = 5 * 1024 * 1024 // 5MB
	if fileHeader.Size > maxImageSize {
		return nil, errs.Body(errs.FormParseError, errors.New("image size exceeds the maximum limit of 5MB"))
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, errs.Body(errs.FormParseError, errors.New("failed to decode the image"))
	}

	if img.Bounds().Dx() < width || img.Bounds().Dy() < height {
		return nil, errs.Body(errs.FormParseError, errors.New("image dimensions should be at least 500x600 pixels"))
	}

	return img, nil
}
