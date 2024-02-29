package application

import (
	"context"
	"fmt"
	"github.com/leetatech/leeta_backend/services/gasrefill/domain"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/mailer"
	"github.com/leetatech/leeta_backend/services/library/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type GasRefillHandler struct {
	idGenerator   library.IDGenerator
	tokenHandler  library.TokenHandler
	logger        *zap.Logger
	EmailClient   mailer.MailerClient
	allRepository library.Repositories
}

type GasRefillApplication interface {
	RequestRefill(ctx context.Context, request domain.GasRefillRequest) (*library.DefaultResponse, error)
}

func NewGasRefillApplication(request library.DefaultApplicationRequest) GasRefillApplication {
	return &GasRefillHandler{
		idGenerator:   library.NewIDGenerator(),
		logger:        request.Logger,
		tokenHandler:  request.TokenHandler,
		EmailClient:   request.EmailClient,
		allRepository: request.AllRepository,
	}
}

func (r GasRefillHandler) RequestRefill(ctx context.Context, request domain.GasRefillRequest) (*library.DefaultResponse, error) {

	claims, err := r.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	if !request.Guest && request.CustomerID != "" {
		customer, err := r.allRepository.UserRepository.GetCustomerByID(claims.UserID)
		if err != nil {
			r.logger.Error("error getting customer", zap.Error(err))
			return nil, err
		}

		if request.ShippingInfo.ForMe {
			request.ShippingInfo = r.forMeCheck(request.ShippingInfo, fmt.Sprintf("%s %s", customer.FirstName, customer.LastName), customer.Phone.Number, customer.Email.Address)
		}
	}

	if request.Guest && request.GuestBioData.Email != "" {
		request.CustomerID = claims.SessionID
		request.GuestBioData.SessionID = claims.SessionID
		if request.ShippingInfo.ForMe {
			request.ShippingInfo = r.forMeCheck(request.ShippingInfo, fmt.Sprintf("%s %s", request.GuestBioData.FirstName, request.GuestBioData.LastName), request.GuestBioData.Phone, request.GuestBioData.Email)
		}
	}

	vendorID, err := r.requestRefill(ctx, request)
	if err != nil {
		return nil, err
	}

	cartItems := models.CartItem{
		ID:         r.idGenerator.Generate(),
		CustomerID: request.CustomerID,
		SessionID:  claims.SessionID,
		ProductID:  request.RefillDetails.ProductID,
		VendorID:   vendorID,
		Weight:     request.RefillDetails.Weight,
		AmountPaid: request.RefillDetails.AmountPaid,
	}

	cart, err := r.allRepository.CartRepository.GetCartBySessionOrCustomerID(ctx, request.CustomerID)
	if err != nil {
		r.logger.Error("error getting cart", zap.Error(err))
		switch err {
		case mongo.ErrNoDocuments:
			err = r.allRepository.CartRepository.AddToCart(ctx, models.Cart{
				ID:          r.idGenerator.Generate(),
				CustomerID:  request.CustomerID,
				CartItems:   []models.CartItem{cartItems},
				DeliveryFee: 0,
				Total:       request.RefillDetails.AmountPaid + 0,
				Status:      models.CartActive,
				StatusTs:    time.Now().Unix(),
				Ts:          time.Now().Unix(),
			})
			return &library.DefaultResponse{Success: "success", Message: "Order successfully created"}, nil
		default:
			return nil, err
		}
	}

	fmt.Println("cartsssss", cart)
	cart.CartItems = append(cart.CartItems, cartItems)
	cart.Total += request.RefillDetails.AmountPaid
	cart.StatusTs = time.Now().Unix()
	err = r.allRepository.CartRepository.AddToCartItem(ctx, cart.ID, cartItems, cart.Total, cart.StatusTs)
	if err != nil {
		r.logger.Error("error adding to cart", zap.Error(err))
		return nil, err
	}

	return &library.DefaultResponse{Success: "success", Message: "Order successfully created"}, nil
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

func (r *GasRefillHandler) requestRefill(ctx context.Context, request domain.GasRefillRequest) (string, error) {
	product, err := r.allRepository.ProductRepository.GetProductByID(ctx, request.RefillDetails.ProductID)
	if err != nil {
		r.logger.Error("error getting product", zap.Error(err))
		return "", err
	}

	vendor, err := r.allRepository.UserRepository.GetVendorByID(product.VendorID)
	if err != nil {
		r.logger.Error("error getting vendor", zap.Error(err))
		return "", err
	}

	refill := models.GasRefill{
		ID:           r.idGenerator.Generate(),
		Guest:        request.Guest,
		GuestBioData: request.GuestBioData,
		CustomerID:   request.CustomerID,
		RefillDetails: models.RefillDetails{
			ProductID:  request.RefillDetails.ProductID,
			VendorID:   vendor.ID,
			Weight:     request.RefillDetails.Weight,
			AmountPaid: request.RefillDetails.AmountPaid,
			GasType:    request.RefillDetails.GasType,
		},
		ShippingInfo: request.ShippingInfo,
		Status:       models.RefillPending,
		StatusTs:     time.Now().Unix(),
		Ts:           time.Now().Unix(),
	}

	err = r.allRepository.GasRefillRepository.RequestRefill(ctx, refill)
	if err != nil {
		r.logger.Error("error requesting refill", zap.Error(err))
		return "", err
	}

	return vendor.ID, nil
}
