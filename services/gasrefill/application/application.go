package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/mailer"
	"github.com/leetatech/leeta_backend/services/gasrefill/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type GasRefillHandler struct {
	idGenerator   pkg.IDGenerator
	tokenHandler  pkg.TokenHandler
	logger        *zap.Logger
	EmailClient   mailer.MailerClient
	allRepository pkg.Repositories
}

type GasRefillApplication interface {
	RequestRefill(ctx context.Context, request domain.GasRefillRequest) (*pkg.DefaultResponse, error)
}

func NewGasRefillApplication(request pkg.DefaultApplicationRequest) GasRefillApplication {
	return &GasRefillHandler{
		idGenerator:   pkg.NewIDGenerator(),
		logger:        request.Logger,
		tokenHandler:  request.TokenHandler,
		EmailClient:   request.EmailClient,
		allRepository: request.AllRepository,
	}
}

func (r *GasRefillHandler) RequestRefill(ctx context.Context, request domain.GasRefillRequest) (*pkg.DefaultResponse, error) {

	claims, err := r.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	if !request.Guest {
		customer, err := r.allRepository.UserRepository.GetCustomerByID(claims.UserID)
		if err != nil {
			r.logger.Error("error getting customer", zap.Error(err))
			return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
		}

		if request.ShippingInfo.ForMe {
			request.ShippingInfo = r.forMeCheck(request.ShippingInfo, fmt.Sprintf("%s %s", customer.FirstName, customer.LastName), customer.Phone.Number, customer.Email.Address)
		}
	}

	if request.Guest && request.GuestBioData.Email != "" {
		request, err = r.manageGuestRefillSession(ctx, request, claims)
		if err != nil {
			return nil, err
		}
	}

	cart, err := r.allRepository.CartRepository.GetCartByCustomerID(ctx, claims.UserID)
	if err != nil {
		r.logger.Error("error getting cart", zap.Error(err))
		return nil, leetError.ErrorResponseBody(leetError.DatabaseNoRecordError, err)
	}

	err = r.requestRefill(ctx, claims.UserID, request, cart)
	if err != nil {
		return nil, err
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Order successfully created"}, nil
}

func (r *GasRefillHandler) manageGuestRefillSession(ctx context.Context, request domain.GasRefillRequest, claims *pkg.UserClaims) (domain.GasRefillRequest, error) {
	request.GuestBioData.DeviceID = claims.UserID

	cart, terr := r.allRepository.CartRepository.GetCartByDeviceID(ctx, claims.UserID)
	if terr != nil {
		if !errors.Is(terr, mongo.ErrNoDocuments) {
			return domain.GasRefillRequest{}, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, terr)
		}
	}

	if cart != nil {
		ts := time.Unix(cart.Ts, 0)
		expectedTime := ts.Add(24 * time.Hour)
		if time.Now().After(expectedTime) || cart.CustomerID != claims.UserID {
			err := r.allRepository.CartRepository.InactivateCart(ctx, cart.ID)
			if err != nil {
				r.logger.Error("error inactivating cart", zap.Error(err))
				return domain.GasRefillRequest{}, leetError.ErrorResponseBody(leetError.DatabaseError, err)
			}
			return domain.GasRefillRequest{}, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, errors.New("guest session expired"))
		}
	}

	request.GuestBioData.SessionID = claims.UserID
	if request.ShippingInfo.ForMe {
		request.ShippingInfo = r.forMeCheck(request.ShippingInfo, fmt.Sprintf("%s %s", request.GuestBioData.FirstName, request.GuestBioData.LastName), request.GuestBioData.Phone, request.GuestBioData.Email)
	}

	return request, nil
}

func (r *GasRefillHandler) forMeCheck(shippingInfo models.ShippingInfo, name, phone, email string) models.ShippingInfo {
	if shippingInfo.RecipientName == "" {
		shippingInfo.RecipientName = name
	}
	if shippingInfo.RecipientPhone == "" {
		shippingInfo.RecipientPhone = phone
	}
	if shippingInfo.RecipientEmail == "" {
		shippingInfo.RecipientEmail = email
	}

	return shippingInfo
}

func (r *GasRefillHandler) requestRefill(ctx context.Context, userID string, request domain.GasRefillRequest, cart *models.Cart) error {
	serviceFee, err := r.calculateCartItemTotal(ctx, cart.CartItems)
	if err != nil {
		return err
	}

	var deliveryFee float64
	totalCost := cart.Total + deliveryFee + serviceFee

	if request.AmountPaid != totalCost {
		return leetError.ErrorResponseBody(leetError.AmountPaidError, errors.New("amount paid does not match total cost"))
	}

	refill := models.GasRefill{
		ID:           r.idGenerator.Generate(),
		Guest:        request.Guest,
		GuestBioData: request.GuestBioData,
		CustomerID:   userID,
		RefillDetails: models.RefillDetails{
			OrderItems: cart.CartItems,
			Status:     models.RefillPending,
			StatusTs:   time.Now().Unix(),
			Ts:         time.Now().Unix(),
		},
		ShippingInfo: request.ShippingInfo,
		AmountPaid:   request.AmountPaid,
		DeliveryFee:  deliveryFee,
		ServiceFee:   serviceFee,
		TotalCost:    totalCost,
		Status:       models.RefillPending,
		StatusTs:     time.Now().Unix(),
		Ts:           time.Now().Unix(),
	}

	err = r.allRepository.GasRefillRepository.RequestRefill(ctx, refill)
	if err != nil {
		r.logger.Error("error requesting refill", zap.Error(err))
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	err = r.allRepository.CartRepository.DeleteCart(ctx, cart.ID)
	if err != nil {
		r.logger.Error("error inactivating cart", zap.Error(err))
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}

func (r *GasRefillHandler) calculateCartItemTotal(ctx context.Context, items []models.CartItem) (float64, error) {
	var serviceFee float64

	fees, err := r.allRepository.FeesRepository.GetFees(ctx, models.FeesActive)
	if err != nil {
		r.logger.Error("get fees", zap.Error(err))
		return 0, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	for _, item := range items {
		for _, fee := range fees {
			if fee.ProductID == item.ProductID {
				serviceFee += fee.ServiceFee
			}
		}
	}

	return serviceFee, nil
}
