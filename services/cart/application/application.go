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
	IncreaseCartItemQuantity(ctx context.Context, request domain.UpdateCartItemQuantityRequest) (*pkg.DefaultResponse, error)
	DecreaseCartItemQuantity(ctx context.Context, request domain.UpdateCartItemQuantityRequest) (*pkg.DefaultResponse, error)
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
		return nil, err
	}

	fee, err := c.allRepository.FeesRepository.GetFeeByProductID(ctx, product.ID, models.FeesActive)
	if err != nil {
		return nil, err
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
		return nil, errors.New("invalid product")
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
				return err
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

func (c CartAppHandler) IncreaseCartItemQuantity(ctx context.Context, request domain.UpdateCartItemQuantityRequest) (*pkg.DefaultResponse, error) {
	var updateRequest domain.UpdateCartItemQuantity
	fee, err := c.updateCartItemQuantity(ctx, request)
	if err != nil {
		return nil, err
	}

	if request.Quantity != 0 {
		updateRequest.CartItemID = request.CartItemID
		updateRequest.Quantity = request.Quantity
		updateRequest.ItemTotalCost = float64(request.Quantity) * fee.CostPerQty
	} else {
		updateRequest.CartItemID = request.CartItemID
		updateRequest.Weight = request.Weight
		updateRequest.ItemTotalCost = request.Weight * fee.CostPerKg
	}

	err = c.allRepository.CartRepository.UpdateCartItemQuantity(ctx, updateRequest)
	if err != nil {
		c.logger.Error("error updating cart item quantity", zap.Error(err))
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Successfully updated cart item quantity"}, nil
}

func (c CartAppHandler) DecreaseCartItemQuantity(ctx context.Context, request domain.UpdateCartItemQuantityRequest) (*pkg.DefaultResponse, error) {
	var updateRequest domain.UpdateCartItemQuantity

	fee, err := c.updateCartItemQuantity(ctx, request)
	if err != nil {
		return nil, err
	}

	if request.Quantity != 0 {
		updateRequest.CartItemID = request.CartItemID
		updateRequest.Quantity = -request.Quantity
		itemTotal := float64(request.Quantity) * fee.CostPerQty
		updateRequest.ItemTotalCost = -itemTotal
	} else {
		updateRequest.CartItemID = request.CartItemID
		updateRequest.Weight = -request.Weight
		itemTotal := request.Weight * fee.CostPerKg
		updateRequest.ItemTotalCost = -itemTotal
	}

	err = c.allRepository.CartRepository.UpdateCartItemQuantity(ctx, updateRequest)
	if err != nil {
		c.logger.Error("error updating cart item quantity", zap.Error(err))
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Successfully updated cart item quantity"}, nil
}

func (c CartAppHandler) updateCartItemQuantity(ctx context.Context, request domain.UpdateCartItemQuantityRequest) (*models.Fee, error) {
	_, err := c.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	var (
		productID        string
		quantityErrorMsg = "invalid quantity"
	)

	cart, err := c.allRepository.CartRepository.GetCartByCartItemID(ctx, request.CartItemID)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	if request.Quantity == 0 && request.Weight == 0 {
		return nil, leetError.ErrorResponseBody(leetError.CartItemRequestQuantityError, errors.New(quantityErrorMsg))
	}

	for _, item := range cart.CartItems {
		if item.ID == request.CartItemID {
			productID = item.ProductID
			if item.Quantity == 0 && item.Weight == 0 {
				return nil, leetError.ErrorResponseBody(leetError.CartItemQuantityError, errors.New(quantityErrorMsg))
			}
		}
	}

	fee, err := c.allRepository.FeesRepository.GetFeeByProductID(ctx, productID, models.FeesActive)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return fee, nil
}
