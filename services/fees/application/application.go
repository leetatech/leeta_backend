package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/paging"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/pkg/idgenerator"
	"github.com/leetatech/leeta_backend/pkg/jwtmiddleware"
	"github.com/leetatech/leeta_backend/pkg/mailer/aws"
	"github.com/leetatech/leeta_backend/services/fees/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type FeesManager struct {
	idgenerator       idgenerator.Generator
	jwtManager        jwtmiddleware.Manager
	EmailClient       aws.MailClient
	repositoryManager pkg.RepositoryManager
}

type Fees interface {
	HandleFeeQuotationRequest(ctx context.Context, request domain.FeeQuotationRequest) (*pkg.DefaultResponse, error)
	Fees(ctx context.Context, request query.ResultSelector) ([]models.Fee, uint64, error)
}

func New(applicationContext pkg.ApplicationContext) Fees {
	return &FeesManager{
		idgenerator:       idgenerator.New(),
		jwtManager:        applicationContext.JwtManager,
		EmailClient:       applicationContext.MailClient,
		repositoryManager: applicationContext.RepositoryManager,
	}
}

func (fm *FeesManager) HandleFeeQuotationRequest(ctx context.Context, request domain.FeeQuotationRequest) (*pkg.DefaultResponse, error) {
	lga, err := fm.typeValidation(ctx, request)
	if err != nil {
		return nil, err
	}

	newFees := models.Fee{
		ID:        fm.idgenerator.Generate(),
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

	fees, _, err := fm.repositoryManager.FeesRepository.Fees(ctx, getRequest)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = fm.repositoryManager.FeesRepository.Create(ctx, newFees)
			if err != nil {
				return nil, errs.Body(errs.DatabaseError, fmt.Errorf("error creating fees: %w", err))
			}

			return &pkg.DefaultResponse{Success: "success", Message: "Fee created successfully"}, nil
		}
		return nil, errs.Body(errs.DatabaseError, err)
	}

	if fees != nil {
		err = fm.repositoryManager.FeesRepository.Update(ctx, models.FeesInactive, request.FeeType, *lga, request.ProductID)
		if err != nil {
			return nil, errs.Body(errs.DatabaseError, fmt.Errorf("error updating fees: %w", err))
		}
		err = fm.repositoryManager.FeesRepository.Create(ctx, newFees)
		if err != nil {
			return nil, errs.Body(errs.DatabaseError, fmt.Errorf("create fees after inactivating stale fee: %w", err))
		}
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Fee created successfully"}, nil
}

func (fm *FeesManager) validateProductFeeRequest(ctx context.Context, request domain.FeeQuotationRequest) error {
	if request.FeeType == models.ProductFee && request.ProductID != "" {
		product, err := fm.repositoryManager.ProductRepository.Product(ctx, request.ProductID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return errs.Body(errs.InvalidProductIdError, fmt.Errorf("no product found with product id %s: %w", request.ProductID, err))
			}
			return errs.Body(errs.DatabaseError, err)
		}

		switch product.ParentCategory {
		case models.LNGProductCategory, models.LPGProductCategory:
			if request.Cost.CostPerKG <= 0 {
				return errs.Body(errs.InvalidRequestError, errors.New("cost per kg is required for product fee"))
			}

		default:
			if request.Cost.CostPerQt <= 0 {
				return errs.Body(errs.InvalidRequestError, errors.New("cost per quantity is required for product fee"))
			}
		}
	}

	return nil
}

func (fm *FeesManager) validateLGA(ctx context.Context, lga models.LGA) error {

	state, err := fm.repositoryManager.StatesRepository.GetState(ctx, lga.State)
	if err != nil {
		return errs.Body(errs.DatabaseError, err)
	}

	// check if the lga exists in the state
	for _, eachLGA := range state.Lgas {
		if eachLGA == lga.LGA {
			return nil
		}
	}

	return errs.Body(errs.InvalidRequestError, errors.New("invalid lga"))
}

func (fm *FeesManager) typeValidation(ctx context.Context, request domain.FeeQuotationRequest) (*models.LGA, error) {

	err := fm.validateProductFeeRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	var lga models.LGA

	if request.FeeType == models.DeliveryFee && request.LGA.State != "" && request.LGA.LGA != "" {
		lga = models.LGA{LGA: request.LGA.LGA, State: strings.ToUpper(request.LGA.State)}
		err = fm.validateLGA(ctx, lga)
		if err != nil {
			return nil, err
		}
	}

	return &lga, nil
}

func (fm *FeesManager) Fees(ctx context.Context, request query.ResultSelector) ([]models.Fee, uint64, error) {
	fees, totalRecord, err := fm.repositoryManager.FeesRepository.Fees(ctx, request)
	if err != nil {
		return nil, 0, errs.Body(errs.InternalError, fmt.Errorf("error fetching fees: %w", err))
	}

	for _, field := range request.Filter.Fields {
		if field.Name == "lga" && len(fees) == 0 {
			return nil, 0, errs.Body(errs.LGANotFoundError, fmt.Errorf("lga not found, leeta is not available in this region: %s", field.Name))
		}
	}

	return fees, totalRecord, nil
}
