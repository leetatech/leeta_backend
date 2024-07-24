package domain

import (
	"context"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/leetatech/leeta_backend/services/models"
)

type FeesRepository interface {
	Create(ctx context.Context, request models.Fee) error
	ByProductID(ctx context.Context, productID string, status models.FeesStatuses) (*models.Fee, error)
	FeesByStatus(ctx context.Context, status models.FeesStatuses) ([]models.Fee, error)
	Update(ctx context.Context, status models.FeesStatuses, feeType models.FeeType, lga models.LGA, productID string) error
	Fees(ctx context.Context, request query.ResultSelector) ([]models.Fee, uint64, error)
}
