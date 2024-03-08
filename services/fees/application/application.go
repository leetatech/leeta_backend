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
	newFees := models.Fee{
		ID:         f.idGenerator.Generate(),
		ProductID:  request.ProductID,
		CostPerKg:  request.CostPerKg,
		CostPerQty: request.CostPerQty,
		ServiceFee: request.ServiceFee,
		Status:     models.CartActive,
		StatusTs:   time.Now().Unix(),
		Ts:         time.Now().Unix(),
	}

	fees, err := f.allRepository.FeesRepository.GetFees(ctx, models.FeesActive)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = f.allRepository.FeesRepository.CreateFees(ctx, newFees)
			if err != nil {
				f.logger.Error("create fees", zap.Error(err))
				return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
			}

			return &library.DefaultResponse{Success: "success", Message: "Fee created successfully"}, nil
		}
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	if fees != nil {
		err = f.allRepository.FeesRepository.UpdateFees(ctx, models.FeesInactive)
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

	return &library.DefaultResponse{Success: "success", Message: "Fee created successfully"}, nil
}

func (f *FeesHandler) GetFees(ctx context.Context) ([]models.Fee, error) {
	fee, err := f.allRepository.FeesRepository.GetFees(ctx, models.FeesActive)
	if err != nil {
		f.logger.Error("get fees", zap.Error(err))
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return fee, nil
}
