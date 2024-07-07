package application

import (
	"context"
	"errors"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/messaging/mailer/awsEmail"
	"github.com/leetatech/leeta_backend/pkg/messaging/mailer/postmarkClient"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/order/domain"
	"go.uber.org/zap"
	"time"
)

type orderAppHandler struct {
	tokenHandler  pkg.TokenHandler
	encryptor     pkg.EncryptorManager
	idGenerator   pkg.IDGenerator
	otpGenerator  pkg.OtpGenerator
	logger        *zap.Logger
	EmailClient   postmarkClient.MailerClient
	AWSClient     awsEmail.AWSEmailClient
	allRepository pkg.Repositories
}

type OrderApplication interface {
	UpdateOrderStatus(ctx context.Context, request domain.UpdateOrderStatusRequest) (*pkg.DefaultResponse, error)
	GetOrderByID(ctx context.Context, id string) (*models.Order, error)
	GetCustomerOrdersByStatus(ctx context.Context, request domain.GetCustomerOrdersRequest) ([]domain.OrderResponse, error)
	ListOrders(ctx context.Context, request query.ResultSelector) ([]models.Order, uint64, error)
	ListOrderStatusHistory(ctx context.Context, orderId string) ([]models.StatusHistory, error)
}

func NewOrderApplication(request pkg.DefaultApplicationRequest) OrderApplication {
	return &orderAppHandler{
		tokenHandler:  request.TokenHandler,
		encryptor:     pkg.NewEncryptor(),
		idGenerator:   pkg.NewIDGenerator(),
		otpGenerator:  pkg.NewOTPGenerator(),
		logger:        request.Logger,
		EmailClient:   request.EmailClient,
		AWSClient:     request.AWSEmailClient,
		allRepository: request.AllRepository,
	}
}

func (o *orderAppHandler) UpdateOrderStatus(ctx context.Context, request domain.UpdateOrderStatusRequest) (*pkg.DefaultResponse, error) {
	claims, err := o.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	if request.Reason == "" {
		return nil, leetError.ErrorResponseBody(leetError.InvalidRequestError, errors.New("reason is required"))
	}

	persistUpdate := domain.PersistOrderUpdate{
		UpdateOrderStatusRequest: request,
		StatusHistory: models.StatusHistory{
			Status:   request.OrderStatus,
			Reason:   request.Reason,
			StatusTs: time.Now().Unix(),
		},
	}

	status, err := models.SetOrderStatus(request.OrderStatus)
	if err != nil {
		return nil, err
	}

	if claims.Role == models.VendorCategory || claims.Role == models.AdminCategory {
		err = o.allRepository.OrderRepository.UpdateOrderStatus(ctx, persistUpdate)
		if err != nil {
			return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
		}

		return &pkg.DefaultResponse{Success: "success", Message: "Order status updated successfully"}, nil
	}

	if status != models.OrderCancelled {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, errors.New("you cannot update this status"))
	}

	err = o.allRepository.OrderRepository.UpdateOrderStatus(ctx, persistUpdate)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Order status updated successfully"}, nil
}

func (o *orderAppHandler) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	_, err := o.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	order, err := o.allRepository.OrderRepository.GetOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *orderAppHandler) GetCustomerOrdersByStatus(ctx context.Context, request domain.GetCustomerOrdersRequest) ([]domain.OrderResponse, error) {
	claims, err := o.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	orders, err := o.allRepository.OrderRepository.GetCustomerOrdersByStatus(ctx, domain.GetCustomerOrders{UserId: claims.UserID, GetCustomerOrdersRequest: request})
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderAppHandler) ListOrders(ctx context.Context, request query.ResultSelector) ([]models.Order, uint64, error) {
	claims, err := o.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, 0, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	orders, totalRecord, err := o.allRepository.OrderRepository.ListOrders(ctx, request, claims.UserID)
	if err != nil {
		return nil, 0, err
	}

	return orders, totalRecord, nil
}

func (o *orderAppHandler) ListOrderStatusHistory(ctx context.Context, orderId string) ([]models.StatusHistory, error) {
	_, err := o.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	statusHistory, err := o.allRepository.OrderRepository.ListOrderStatusHistory(ctx, orderId)
	if err != nil {
		return nil, err
	}

	return statusHistory, nil
}
