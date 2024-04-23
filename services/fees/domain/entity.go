package domain

import (
	"errors"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/services/models"
)

type FeeQuotationRequest struct {
	Cost      models.Cost    `json:"cost" bson:"cost"`
	FeeType   models.FeeType `json:"fee_type" bson:"fee_type"`
	LGA       models.LGA     `json:"lga,omitempty" bson:"lga"`
	ProductID string         `json:"product_id,omitempty" bson:"product_id"`
} // @name FeeQuotationRequest

type GetTypedFeesRequest struct {
	FeeTypes    []models.FeeType      `json:"fee_types"`
	FeeStatuses []models.FeesStatuses `json:"fee_statuses"`
	LGA         *models.LGA           `json:"lga"`
	ProductID   *string               `json:"product_id"`
} //@name GetTypedFeesRequest

func (request FeeQuotationRequest) FeeTypeValidation() (FeeQuotationRequest, error) {
	switch request.FeeType {
	case models.ServiceFee:
		request.LGA = models.LGA{}
		request.ProductID = ""
		request.Cost.CostPerQt = 0
		request.Cost.CostPerKG = 0
		if request.Cost.CostPerType == 0 {
			return FeeQuotationRequest{}, leetError.ErrorResponseBody(leetError.InvalidRequestError, errors.New("cost per type is required for service fee"))
		}
	case models.ProductFee:
		request.LGA = models.LGA{}
		request.Cost.CostPerType = 0
		if request.ProductID == "" {
			return FeeQuotationRequest{}, leetError.ErrorResponseBody(leetError.InvalidRequestError, errors.New("product id, cost per kg and cost per quantity is required for product fee"))
		}

	case models.DeliveryFee:
		request.ProductID = ""
		request.Cost.CostPerQt = 0
		request.Cost.CostPerKG = 0
		if request.LGA.LGA == "" || request.Cost.CostPerType == 0 {
			return FeeQuotationRequest{}, leetError.ErrorResponseBody(leetError.InvalidRequestError, errors.New("lga and cost per type is required for delivery fee"))
		}
	}

	return request, nil
}
