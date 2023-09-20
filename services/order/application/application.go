package application

import (
	"context"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/mailer"
	"github.com/leetatech/leeta_backend/services/library/models"
	"github.com/leetatech/leeta_backend/services/order/domain"
	"go.uber.org/zap"
	"time"
)

type orderAppHandler struct {
	tokenHandler  library.TokenHandler
	encryptor     library.EncryptorManager
	idGenerator   library.IDGenerator
	otpGenerator  library.OtpGenerator
	logger        *zap.Logger
	EmailClient   mailer.MailerClient
	allRepository library.Repositories
}

type OrderApplication interface {
	CreateOrder(ctx context.Context, request domain.OrderRequest) (*library.DefaultResponse, error)
	UpdateOrderStatus(ctx context.Context, request domain.UpdateOrderStatusRequest) (*library.DefaultResponse, error)
	GetOrderByID(ctx context.Context, id string) (*models.Order, error)
	GetCustomerOrdersByStatus(ctx context.Context, request domain.GetCustomerOrdersRequest) ([]domain.OrderResponse, error)
}

func NewOrderApplication(request library.DefaultApplicationRequest) OrderApplication {
	return &orderAppHandler{
		tokenHandler:  request.TokenHandler,
		encryptor:     library.NewEncryptor(),
		idGenerator:   library.NewIDGenerator(),
		otpGenerator:  library.NewOTPGenerator(),
		logger:        request.Logger,
		EmailClient:   request.EmailClient,
		allRepository: request.AllRepository,
	}
}

func (o orderAppHandler) CreateOrder(ctx context.Context, request domain.OrderRequest) (*library.DefaultResponse, error) {
	claims, err := o.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	_, err = o.allRepository.AuthRepository.GetCustomerByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}

	product, err := o.allRepository.ProductRepository.GetProductByID(ctx, request.ProductID)
	if err != nil {
		return nil, err
	}

	vendor, err := o.allRepository.UserRepository.GetVendorByID(product.VendorID)
	if err != nil {
		return nil, err
	}

	//TODO delivery fees based on location
	deliveryFee := 1000.00
	totalCost := deliveryFee + product.Vat

	order := models.Order{
		ID:          o.idGenerator.Generate(),
		ProductID:   request.ProductID,
		CustomerID:  claims.UserID,
		VendorID:    vendor.ID,
		VAT:         product.Vat,
		DeliveryFee: deliveryFee,
		Total:       totalCost,
		Status:      models.OrderPending,
		StatusTs:    time.Now().Unix(),
		Ts:          time.Now().Unix(),
	}
	err = o.allRepository.OrderRepository.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return &library.DefaultResponse{Success: "success", Message: "Order successfully created"}, nil
}

func (o orderAppHandler) UpdateOrderStatus(ctx context.Context, request domain.UpdateOrderStatusRequest) (*library.DefaultResponse, error) {
	claims, err := o.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	if claims.Role != models.AdminCategory {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}
	_, err = o.allRepository.AuthRepository.GetAdminByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}

	_, err = models.SetOrderStatus(request.OrderStatus)
	if err != nil {
		return nil, err
	}

	err = o.allRepository.OrderRepository.UpdateOrderStatus(ctx, request)
	if err != nil {
		return nil, err
	}

	return &library.DefaultResponse{Success: "success", Message: "Order status successfully updated"}, nil
}

func (o orderAppHandler) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
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

func (o orderAppHandler) GetCustomerOrdersByStatus(ctx context.Context, request domain.GetCustomerOrdersRequest) ([]domain.OrderResponse, error) {
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
