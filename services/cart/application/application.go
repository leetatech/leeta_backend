package application

import (
	"context"
	"errors"
	"github.com/leetatech/leeta_backend/services/cart/domain"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/mailer"
	"github.com/leetatech/leeta_backend/services/library/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type CartAppHandler struct {
	idGenerator   library.IDGenerator
	tokenHandler  library.TokenHandler
	logger        *zap.Logger
	EmailClient   mailer.MailerClient
	allRepository library.Repositories
}

type CartApplication interface {
	InactivateCart(ctx context.Context, request domain.InactivateCart) (*library.DefaultResponse, error)
	AddToCart(ctx context.Context, request domain.AddToCartRequest) (*library.DefaultResponse, error)
}

func NewCartApplication(request library.DefaultApplicationRequest) CartApplication {
	return &CartAppHandler{
		idGenerator:   library.NewIDGenerator(),
		logger:        request.Logger,
		tokenHandler:  request.TokenHandler,
		EmailClient:   request.EmailClient,
		allRepository: request.AllRepository,
	}
}

func (c CartAppHandler) InactivateCart(ctx context.Context, request domain.InactivateCart) (*library.DefaultResponse, error) {
	err := c.allRepository.CartRepository.InactivateCart(ctx, request.ID)
	if err != nil {
		return nil, err
	}
	return &library.DefaultResponse{Success: "success", Message: "Cart inactivated successfully"}, nil
}

func (c CartAppHandler) AddToCart(ctx context.Context, request domain.AddToCartRequest) (*library.DefaultResponse, error) {
	claims, err := c.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	var deviceID string
	if request.Guest {
		deviceID = claims.UserID

		err := c.manageGuestCartSession(ctx, deviceID, claims)
		if err != nil {
			return nil, err
		}
		claims.UserID = claims.SessionID
	}

	product, err := c.allRepository.ProductRepository.GetProductByID(ctx, request.RefillDetails.ProductID)
	if err != nil {
		c.logger.Error("error getting product", zap.Error(err))
		return nil, err
	}

	fee, err := c.allRepository.FeesRepository.GetFeeByProductID(ctx, product.ID, models.CartActive)
	if err != nil {
		return nil, err
	}
	cartItems := models.CartItem{
		ID:        c.idGenerator.Generate(),
		ProductID: request.RefillDetails.ProductID,
		VendorID:  product.VendorID,
		Weight:    request.RefillDetails.Weight,
		TotalCost: fee.CostPerKg * float64(request.RefillDetails.Weight),
	}

	cart, err := c.allRepository.CartRepository.GetCartByCustomerID(ctx, claims.UserID)
	if err != nil {
		c.logger.Error("error getting cart", zap.Error(err))
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			err = c.allRepository.CartRepository.AddToCart(ctx, models.Cart{
				ID:         c.idGenerator.Generate(),
				CustomerID: claims.UserID,
				DeviceID:   deviceID,
				CartItems:  []models.CartItem{cartItems},
				Total:      cartItems.TotalCost,
				Status:     models.CartActive,
				StatusTs:   time.Now().Unix(),
				Ts:         time.Now().Unix(),
			})
			return &library.DefaultResponse{Success: "success", Message: "Successfully added item to cart"}, nil
		default:
			return nil, err
		}
	}

	cart.CartItems = append(cart.CartItems, cartItems)

	cart.Total, err = c.calculateCartItemTotal(ctx, cart.CartItems)
	if err != nil {
		c.logger.Error("error calculating cart total", zap.Error(err))
		return nil, err
	}

	cart.StatusTs = time.Now().Unix()
	err = c.allRepository.CartRepository.AddToCartItem(ctx, cart.ID, cartItems, cart.Total, cart.StatusTs)
	if err != nil {
		c.logger.Error("error adding to cart", zap.Error(err))
		return nil, err
	}

	return &library.DefaultResponse{Success: "success", Message: "Successfully added item to cart"}, nil
}

func (c CartAppHandler) manageGuestCartSession(ctx context.Context, deviceID string, claims *library.UserClaims) error {
	cart, terr := c.allRepository.CartRepository.GetCartByDeviceID(ctx, deviceID)
	if terr != nil {
		if !errors.Is(terr, mongo.ErrNoDocuments) {
			return leetError.ErrorResponseBody(leetError.ErrorUnauthorized, terr)
		}
	}

	if cart != nil {
		ts := time.Unix(cart.Ts, 0)
		expectedTime := ts.Add(24 * time.Hour)
		if time.Now().After(expectedTime) || cart.CustomerID != claims.SessionID {
			err := c.allRepository.CartRepository.InactivateCart(ctx, cart.ID)
			if err != nil {
				return err
			}
			return leetError.ErrorResponseBody(leetError.ErrorUnauthorized, errors.New("guest session expired"))
		}
	}

	return nil
}

func (c CartAppHandler) calculateCartItemTotal(ctx context.Context, items []models.CartItem) (float64, error) {
	var total float64

	fees, err := c.allRepository.FeesRepository.GetFees(ctx, models.CartActive)
	if err != nil {
		return 0, err
	}

	for _, item := range items {
		for _, fee := range fees {
			if fee.ProductID == item.ProductID {
				total += fee.CostPerKg * float64(item.Weight)
			}
		}
	}

	return total, nil
}
