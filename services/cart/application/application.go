package application

import (
	"context"
	"errors"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/mailer"
	"github.com/leetatech/leeta_backend/services/cart/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type CartAppHandler struct {
	idGenerator   pkg.IDGenerator
	tokenHandler  pkg.TokenHandler
	logger        *zap.Logger
	EmailClient   mailer.MailerClient
	allRepository pkg.Repositories
}

type CartApplication interface {
	InactivateCart(ctx context.Context, request domain.InactivateCart) (*pkg.DefaultResponse, error)
	AddToCart(ctx context.Context, request domain.AddToCartRequest) (*pkg.DefaultResponse, error)
	DeleteCartItem(ctx context.Context, request domain.DeleteCartItemRequest) (*pkg.DefaultResponse, error)
}

func NewCartApplication(request pkg.DefaultApplicationRequest) CartApplication {
	return &CartAppHandler{
		idGenerator:   pkg.NewIDGenerator(),
		logger:        request.Logger,
		tokenHandler:  request.TokenHandler,
		EmailClient:   request.EmailClient,
		allRepository: request.AllRepository,
	}
}

func (c CartAppHandler) InactivateCart(ctx context.Context, request domain.InactivateCart) (*pkg.DefaultResponse, error) {
	err := c.allRepository.CartRepository.InactivateCart(ctx, request.ID)
	if err != nil {
		return nil, err
	}
	return &pkg.DefaultResponse{Success: "success", Message: "Cart inactivated successfully"}, nil
}

func (c CartAppHandler) AddToCart(ctx context.Context, request domain.AddToCartRequest) (*pkg.DefaultResponse, error) {
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

	product, err := c.allRepository.ProductRepository.GetProductByID(ctx, request.CartDetails.ProductID)
	if err != nil {
		c.logger.Error("error getting product", zap.Error(err))
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	fee, err := c.allRepository.FeesRepository.GetFeeByProductID(ctx, product.ID, models.FeesActive)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	cartItems := models.CartItem{
		ID:        c.idGenerator.Generate(),
		ProductID: request.CartDetails.ProductID,
		VendorID:  product.VendorID,
		Weight:    request.CartDetails.Weight,
	}

	cartItems.TotalCost = cartItems.CalculateCartFee(fee)
	if cartItems.TotalCost == 0 {
		c.logger.Error("invalid product id")
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, errors.New("invalid product id"))
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
			return &pkg.DefaultResponse{Success: "success", Message: "Successfully added item to cart"}, nil
		default:
			return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
		}
	}

	cart.CartItems = append(cart.CartItems, cartItems)

	cart.Total, err = c.calculateCartItemTotal(ctx, cart.CartItems)
	if err != nil {
		c.logger.Error("error calculating cart total", zap.Error(err))
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	cart.StatusTs = time.Now().Unix()
	err = c.allRepository.CartRepository.AddToCartItem(ctx, cart.ID, cartItems, cart.Total, cart.StatusTs)
	if err != nil {
		c.logger.Error("error adding to cart", zap.Error(err))
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Successfully added item to cart"}, nil
}

func (c CartAppHandler) manageGuestCartSession(ctx context.Context, deviceID string, claims *pkg.UserClaims) error {
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
				c.logger.Error("error inactivating cart", zap.Error(err))
				return leetError.ErrorResponseBody(leetError.DatabaseError, err)
			}
			return leetError.ErrorResponseBody(leetError.ErrorUnauthorized, errors.New("guest session expired"))
		}
	}

	return nil
}

func (c CartAppHandler) calculateCartItemTotal(ctx context.Context, items []models.CartItem) (float64, error) {
	var total float64

	fees, err := c.allRepository.FeesRepository.GetFees(ctx, models.FeesActive)
	if err != nil {
		return 0, err
	}

	for _, item := range items {
		for _, fee := range fees {
			if fee.ProductID == item.ProductID {
				total += item.CalculateCartFee(&fee)
			}
		}
	}

	return total, nil
}

func (c CartAppHandler) DeleteCartItem(ctx context.Context, request domain.DeleteCartItemRequest) (*pkg.DefaultResponse, error) {
	_, err := c.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	cart, err := c.allRepository.CartRepository.GetCartByCartItemID(ctx, request.CartItemID)
	if err != nil {
		c.logger.Error("error getting cart by cart item id", zap.Error(err))
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	var itemTotalCost float64
	for _, item := range cart.CartItems {
		if item.ID == request.CartItemID {
			itemTotalCost = item.TotalCost
		}
	}

	err = c.allRepository.CartRepository.DeleteCartItem(ctx, request.CartItemID, itemTotalCost)
	if err != nil {
		c.logger.Error("error deleting cart item", zap.Error(err))
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Successfully deleted item from cart"}, nil
}
