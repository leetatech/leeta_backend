package application

import (
	"context"
	"errors"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/messaging/mailer/awsEmail"
	"github.com/leetatech/leeta_backend/pkg/messaging/mailer/postmarkClient"
	"github.com/leetatech/leeta_backend/pkg/encrypto"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/pkg/idgenerator"
	"github.com/leetatech/leeta_backend/pkg/jwtmiddleware"
	"github.com/leetatech/leeta_backend/pkg/mailer/aws"
	"github.com/leetatech/leeta_backend/pkg/otp"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/order/domain"
	"time"
)

type orderAppHandler struct {
	jwtManager    jwtmiddleware.Manager
	encryptor     encrypto.Manager
	idGenerator   idgenerator.Generator
	otpGenerator  otp.Generator
	EmailClient   aws.MailClient
	allRepository pkg.RepositoryManager
}

type Order interface {
	UpdateOrderStatus(ctx context.Context, request domain.UpdateStatusRequest) (*pkg.DefaultResponse, error)
	GetOrderByID(ctx context.Context, id string) (*models.Order, error)
	GetCustomerOrdersByStatus(ctx context.Context, request domain.GetCustomerOrdersRequest) ([]domain.Response, error)
	ListOrders(ctx context.Context, request query.ResultSelector) ([]models.Order, uint64, error)
	ListOrderStatusHistory(ctx context.Context, orderId string) ([]models.StatusHistory, error)
}

func New(request pkg.ApplicationContext) Order {
	return &orderAppHandler{
		jwtManager:    request.JwtManager,
		encryptor:     encrypto.New(),
		idGenerator:   idgenerator.New(),
		otpGenerator:  otp.New(),
		EmailClient:   request.MailClient,
		allRepository: request.RepositoryManager,
	}
}

func (o *orderAppHandler) UpdateOrderStatus(ctx context.Context, request domain.UpdateStatusRequest) (*pkg.DefaultResponse, error) {
	claims, err := o.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return nil, errs.Body(errs.ErrorUnauthorized, err)
	}

	if request.Reason == "" {
		return nil, errs.Body(errs.InvalidRequestError, errors.New("reason is required"))
	}

	persistUpdate := domain.PersistOrderUpdate{
		UpdateStatusRequest: request,
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
		err = o.allRepository.OrderRepository.UpdateStatus(ctx, persistUpdate)
		if err != nil {
			return nil, errs.Body(errs.DatabaseError, err)
		}

		return &pkg.DefaultResponse{Success: "success", Message: "Order status updated successfully"}, nil
	}

	if status != models.OrderCancelled {
		return nil, errs.Body(errs.ErrorUnauthorized, errors.New("you cannot update this status"))
	}

	err = o.allRepository.OrderRepository.UpdateStatus(ctx, persistUpdate)
	if err != nil {
		return nil, errs.Body(errs.DatabaseError, err)
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Order status updated successfully"}, nil
}

func (o *orderAppHandler) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	_, err := o.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return nil, errs.Body(errs.ErrorUnauthorized, err)
	}

	order, err := o.allRepository.OrderRepository.OrderByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *orderAppHandler) GetCustomerOrdersByStatus(ctx context.Context, request domain.GetCustomerOrdersRequest) ([]domain.Response, error) {
	claims, err := o.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return nil, errs.Body(errs.ErrorUnauthorized, err)
	}

	orders, err := o.allRepository.OrderRepository.OrdersByStatus(ctx, domain.GetCustomerOrders{UserId: claims.UserID, GetCustomerOrdersRequest: request})
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderAppHandler) ListOrders(ctx context.Context, request query.ResultSelector) ([]models.Order, uint64, error) {
	claims, err := o.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return nil, 0, errs.Body(errs.ErrorUnauthorized, err)
	}

	orders, totalRecord, err := o.allRepository.OrderRepository.Orders(ctx, request, claims.UserID)
	if err != nil {
		return nil, 0, err
	}

	return orders, totalRecord, nil
}

func (o *orderAppHandler) ListOrderStatusHistory(ctx context.Context, orderId string) ([]models.StatusHistory, error) {
	_, err := o.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return nil, errs.Body(errs.ErrorUnauthorized, err)
	}

	statusHistory, err := o.allRepository.OrderRepository.OrderStatusHistory(ctx, orderId)
	if err != nil {
		return nil, err
	}

	return statusHistory, nil
}
