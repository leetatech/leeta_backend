package application

import (
	"context"
	"errors"
	"github.com/leetatech/leeta_backend/services/fees/domain"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/mailer"
	"github.com/leetatech/leeta_backend/services/library/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type FeesHandler struct {
	idGenerator   library.IDGenerator
	tokenHandler  library.TokenHandler
	logger        *zap.Logger
	EmailClient   mailer.MailerClient
	allRepository library.Repositories
}

type FeesApplication interface {
	FeeQuotation(ctx context.Context, request domain.FeeQuotationRequest) (*library.DefaultResponse, error)
	GetFees(ctx context.Context) ([]models.Fee, error)
	GetFeeByProductID(ctx context.Context, productID string) (*models.Fee, error)
}

func NewFeesApplication(request library.DefaultApplicationRequest) FeesApplication {
	return &FeesHandler{
		idGenerator:   library.NewIDGenerator(),
		logger:        request.Logger,
		tokenHandler:  request.TokenHandler,
		EmailClient:   request.EmailClient,
		allRepository: request.AllRepository,
	}
}

func (f *FeesHandler) FeeQuotation(ctx context.Context, request domain.FeeQuotationRequest) (*library.DefaultResponse, error) {
	_, err := f.allRepository.ProductRepository.GetProductByID(ctx, request.ProductID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, leetError.ErrorResponseBody(leetError.DatabaseNoRecordError, err)
		}
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	newFees := models.Fee{
		ID:         f.idGenerator.Generate(),
		ProductID:  request.ProductID,
		CostPerKg:  request.CostPerKg,
		CostPerQty: request.CostPerQty,
		Status:     models.CartActive,
		StatusTs:   time.Now().Unix(),
		Ts:         time.Now().Unix(),
	}

	fees, err := f.allRepository.FeesRepository.GetFees(ctx, models.FeesActive)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = f.allRepository.FeesRepository.CreateFees(ctx, newFees)
			if err != nil {
				return nil, err
			}

			return &library.DefaultResponse{Success: "success", Message: "Fee created successfully"}, nil
		}
		return nil, err
	}

	if fees != nil {
		err = f.allRepository.FeesRepository.UpdateFees(ctx, models.FeesInactive)
		if err != nil {
			return nil, err
		}
		err = f.allRepository.FeesRepository.CreateFees(ctx, newFees)
		if err != nil {
			return nil, err
		}
	}

	return &library.DefaultResponse{Success: "success", Message: "Fee created successfully"}, nil
}

func (f *FeesHandler) GetFees(ctx context.Context) ([]models.Fee, error) {
	fee, err := f.allRepository.FeesRepository.GetFees(ctx, models.FeesActive)
	if err != nil {
		return nil, err
	}

	return fee, nil
}

func (f *FeesHandler) GetFeeByProductID(ctx context.Context, productID string) (*models.Fee, error) {
	fee, err := f.allRepository.FeesRepository.GetFeeByProductID(ctx, productID, models.FeesActive)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return fee, nil
}
