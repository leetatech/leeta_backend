package application

import (
	"context"
	"errors"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/paging"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/mailer"
	"github.com/leetatech/leeta_backend/services/fees/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"strings"
	"time"
)

type FeesHandler struct {
	idGenerator   pkg.IDGenerator
	tokenHandler  pkg.TokenHandler
	logger        *zap.Logger
	EmailClient   mailer.MailerClient
	allRepository pkg.Repositories
}

type FeesApplication interface {
	FeeQuotation(ctx context.Context, request domain.FeeQuotationRequest) (*pkg.DefaultResponse, error)
	GetTypedFees(ctx context.Context, request query.ResultSelector) ([]models.Fee, uint64, error)
}

func NewFeesApplication(request pkg.DefaultApplicationRequest) FeesApplication {
	return &FeesHandler{
		idGenerator:   pkg.NewIDGenerator(),
		logger:        request.Logger,
		tokenHandler:  request.TokenHandler,
		EmailClient:   request.EmailClient,
		allRepository: request.AllRepository,
	}
}

func (f *FeesHandler) FeeQuotation(ctx context.Context, request domain.FeeQuotationRequest) (*pkg.DefaultResponse, error) {
	lga, err := f.feeTypeValidation(ctx, request)
	if err != nil {
		return nil, err
	}

	newFees := models.Fee{
		ID:        f.idGenerator.Generate(),
		ProductID: request.ProductID,
		FeeType:   request.FeeType,
		LGA:       *lga,
		Cost: models.Cost{
			CostPerKG:   request.Cost.CostPerKG,
			CostPerQt:   request.Cost.CostPerQt,
			CostPerType: request.Cost.CostPerType,
		},
		Status:   models.FeesActive,
		StatusTs: time.Now().Unix(),
		Ts:       time.Now().Unix(),
	}

	getRequest := query.ResultSelector{
		Filter: &filter.Request{
			Operator: "and",
			Fields: []filter.RequestField{
				{
					Name:  "lga",
					Value: lga,
				},
				{
					Name:  "product_id",
					Value: request.ProductID,
				},
				{
					Name:  "fee_type",
					Value: request.FeeType,
				},
				{
					Name:  "status",
					Value: newFees.Status,
				},
			},
		},
		Paging: &paging.Request{},
	}

	fees, _, err := f.allRepository.FeesRepository.FetchFees(ctx, getRequest)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = f.allRepository.FeesRepository.CreateFees(ctx, newFees)
			if err != nil {
				f.logger.Error("create fees", zap.Error(err))
				return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
			}

			return &pkg.DefaultResponse{Success: "success", Message: "Fee created successfully"}, nil
		}
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	if fees != nil {
		err = f.allRepository.FeesRepository.UpdateFees(ctx, models.FeesInactive, request.FeeType, *lga, request.ProductID)
		if err != nil {
			f.logger.Error("update fees", zap.Error(err))
			return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
		}
		err = f.allRepository.FeesRepository.CreateFees(ctx, newFees)
		if err != nil {
			f.logger.Error("create fees after inactivating the previous one", zap.Error(err))
			return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
		}
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Fee created successfully"}, nil
}

func (f *FeesHandler) validateProductFeeRequest(ctx context.Context, request domain.FeeQuotationRequest) error {
	if request.FeeType == models.ProductFee && request.ProductID != "" {
		product, err := f.allRepository.ProductRepository.GetProductByID(ctx, request.ProductID)
		if err != nil {
			f.logger.Error("get product by id", zap.Error(err))
			if errors.Is(err, mongo.ErrNoDocuments) {
				return leetError.ErrorResponseBody(leetError.InvalidProductIdError, err)
			}
			return leetError.ErrorResponseBody(leetError.DatabaseError, err)
		}

		switch product.ParentCategory {
		case models.LNGProductCategory, models.LPGProductCategory:
			if request.Cost.CostPerKG <= 0 {
				return leetError.ErrorResponseBody(leetError.InvalidRequestError, errors.New("cost per kg is required for product fee"))
			}

		default:
			if request.Cost.CostPerQt <= 0 {
				return leetError.ErrorResponseBody(leetError.InvalidRequestError, errors.New("cost per quantity is required for product fee"))
			}
		}
	}

	return nil
}

func (f *FeesHandler) validateLGA(ctx context.Context, lga models.LGA) error {

	state, err := f.allRepository.StatesRepository.GetState(ctx, lga.State)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	// check if the lga exists in the state
	for _, eachLGA := range state.Lgas {
		if eachLGA == lga.LGA {
			return nil
		}
	}

	return leetError.ErrorResponseBody(leetError.InvalidRequestError, errors.New("invalid lga"))
}

func (f *FeesHandler) feeTypeValidation(ctx context.Context, request domain.FeeQuotationRequest) (*models.LGA, error) {

	err := f.validateProductFeeRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	var lga models.LGA

	if request.FeeType == models.DeliveryFee && request.LGA.State != "" && request.LGA.LGA != "" {
		lga = models.LGA{LGA: request.LGA.LGA, State: strings.ToUpper(request.LGA.State)}
		err = f.validateLGA(ctx, lga)
		if err != nil {
			return nil, err
		}
	}

	return &lga, nil
}

func (f *FeesHandler) GetTypedFees(ctx context.Context, request query.ResultSelector) ([]models.Fee, uint64, error) {
	fees, totalRecord, err := f.allRepository.FeesRepository.FetchFees(ctx, request)
	if err != nil {
		return nil, 0, err
	}

	return fees, totalRecord, nil
}
