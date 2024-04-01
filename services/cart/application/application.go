package application

import (
	"context"
	"errors"
	"fmt"
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
	AddToCart(ctx context.Context, request domain.CartItem) (*pkg.DefaultResponse, error)
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

func (c CartAppHandler) AddToCart(ctx context.Context, request domain.CartItem) (*pkg.DefaultResponse, error) {
	claims, err := c.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting user claims %w", leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err))
	}

	product, err := c.allRepository.ProductRepository.GetProductByID(ctx, request.ProductID)
	if err != nil {
		return nil, fmt.Errorf("error getting product id %s: %w", request.ProductID, err)
	}

	fee, err := c.allRepository.FeesRepository.GetFeeByProductID(ctx, product.ID, models.FeesActive)
	if err != nil {
		return nil, fmt.Errorf("error getting fee from product with id %s: %w", product.ID, err)
	}

	cartItem := models.CartItem{
		ID:        c.idGenerator.Generate(),
		ProductID: request.ProductID,
		VendorID:  product.VendorID,
		Weight:    request.Weight,
		Quantity:  request.Quantity,
	}

	cartItem.TotalCost, err = cartItem.CalculateCartFee(fee)
	if cartItem.TotalCost == 0 {
		return nil, fmt.Errorf("unable to calculate cart fee %w", err)
	}

	cart, err := c.allRepository.CartRepository.GetCartByCustomerID(ctx, claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			addToCartErr := c.allRepository.CartRepository.AddToCart(ctx, models.Cart{
				ID:         c.idGenerator.Generate(),
				CustomerID: claims.SessionID,
				CartItems:  []models.CartItem{cartItem},
				Total:      cartItem.TotalCost,
				Status:     models.CartActive,
				StatusTs:   time.Now().Unix(),
				Ts:         time.Now().Unix(),
			})
			if addToCartErr != nil {
				return nil, fmt.Errorf("error when adding item to cart store %w", err)
			}
			return &pkg.DefaultResponse{Success: "success", Message: "Successfully added item to cart"}, nil
		default:
			return nil, fmt.Errorf("error getting cart item by customer id %w", err)
		}
	}

	cart.CartItems = append(cart.CartItems, cartItem)

	cart.Total, err = c.calculateCartItemTotalCost(ctx, cart.CartItems)
	if err != nil {
		return nil, fmt.Errorf("error calculating cart item total fee %w", err)
	}

	cart.StatusTs = time.Now().Unix()
	err = c.allRepository.CartRepository.AddToCartItem(ctx, cart.ID, cartItem, cart.Total, cart.StatusTs)
	if err != nil {
		return nil, fmt.Errorf("error adding item to cart %w", err)
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Successfully added item to cart"}, nil
}

func (c CartAppHandler) calculateCartItemTotalCost(ctx context.Context, items []models.CartItem) (float64, error) {
	var total float64

	fees, err := c.allRepository.FeesRepository.GetFees(ctx, models.FeesActive)
	if err != nil {
		return 0, err
	}

	for _, item := range items {
		for _, fee := range fees {
			if fee.ProductID == item.ProductID {
				cartTotalFee, err := item.CalculateCartFee(&fee)
				if err != nil {
					return 0, fmt.Errorf("error calculating cart fee %w", err)
				}
				total += cartTotalFee
			}
		}
	}

	return total, nil
}
