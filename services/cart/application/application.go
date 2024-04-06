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
	"math"
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
	UpdateCartItemQuantity(ctx context.Context, request domain.UpdateCartItemQuantityRequest) (*pkg.DefaultResponse, error)
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

	switch product.ParentCategory {
	case models.LNGProductCategory, models.LPGProductCategory:
		if request.Weight == 0 {
			return nil, leetError.ErrorResponseBody(leetError.InvalidRequestError, errors.New("invalid cart item, cart weight cannot be zero"))
		}
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

	cartItem.Cost, err = cartItem.CalculateCartFee(fee)
	if cartItem.Cost == 0 {
		return nil, fmt.Errorf("unable to calculate cart fee %w", err)
	}

	cart, err := c.allRepository.CartRepository.GetCartByCustomerID(ctx, claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			addToCartErr := c.allRepository.CartRepository.AddToCart(ctx, models.Cart{
				ID:         c.idGenerator.Generate(),
				CustomerID: claims.UserID,
				CartItems:  []models.CartItem{cartItem},
				Total:      cartItem.Cost,
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

func (c CartAppHandler) UpdateCartItemQuantity(ctx context.Context, request domain.UpdateCartItemQuantityRequest) (*pkg.DefaultResponse, error) {
	_, err := c.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	resp, updateRequest, err := c.gasProductCategoryAdjustment(ctx, request)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		return resp, nil
	}

	err = c.allRepository.CartRepository.UpdateCartItemQuantity(ctx, *updateRequest)
	if err != nil {
		c.logger.Error("error updating cart item quantity", zap.Error(err))
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Successfully updated cart item quantity"}, nil
}

func (c CartAppHandler) compareStoredAndRequestQuantity(ctx context.Context, request domain.UpdateCartItemQuantityRequest) (*pkg.DefaultResponse, *domain.StoredCartItemDetails, error) {
	var (
		quantityErrorMsg = "invalid quantity"
	)

	cart, err := c.allRepository.CartRepository.GetCartByCartItemID(ctx, request.CartItemID)
	if err != nil {
		return nil, nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	if request.Quantity == 0 {
		return nil, nil, leetError.ErrorResponseBody(leetError.CartItemRequestQuantityError, errors.New(quantityErrorMsg))
	}

	// TODO itemsInCart := len(cart.CartItems)

	var storedItemDetails domain.StoredCartItemDetails

	for _, item := range cart.CartItems {
		if item.ID == request.CartItemID {
			storedItemDetails = domain.StoredCartItemDetails{
				ProductID: item.ProductID,
				Weight:    item.Weight,
				Quantity:  item.Quantity,
			}
			if request.Quantity < 0 {
				if item.Quantity == 1 && request.Quantity < 0 || item.Quantity == int(math.Abs(float64(request.Quantity))) {
					err = c.allRepository.CartRepository.DeleteCartItem(ctx, request.CartItemID, item.Cost)
					if err != nil {
						return nil, nil, err
					}

					// TODO when others are merged if itemsInCart == 1 {}

					return &pkg.DefaultResponse{Success: "success", Message: "Successfully deleted cart item"}, nil, nil
				}
			}
		}
	}

	var product *models.Product
	if storedItemDetails.ProductID != "" {
		product, err = c.allRepository.ProductRepository.GetProductByID(ctx, storedItemDetails.ProductID)
		if err != nil {
			return nil, nil, err
		}
	}

	storedItemDetails.ProductCategory = product.ParentCategory

	return nil, &storedItemDetails, nil
}

func (c CartAppHandler) retrieveProductFee(ctx context.Context, productID string) (*models.Fee, error) {

	fee, err := c.allRepository.FeesRepository.GetFeeByProductID(ctx, productID, models.FeesActive)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return fee, nil
}

func (c CartAppHandler) gasProductCategoryAdjustment(ctx context.Context, request domain.UpdateCartItemQuantityRequest) (*pkg.DefaultResponse, *domain.UpdateCartItemQuantity, error) {
	resp, storedCartItemDetails, err := c.compareStoredAndRequestQuantity(ctx, request)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return resp, nil, nil
	}

	fee, err := c.retrieveProductFee(ctx, storedCartItemDetails.ProductID)
	if err != nil {
		return nil, nil, err
	}

	var updateRequest domain.UpdateCartItemQuantity

	switch storedCartItemDetails.ProductCategory {
	case models.LNGProductCategory, models.LPGProductCategory:
		newQuantity := storedCartItemDetails.Quantity + request.Quantity
		oldItemCost := storedCartItemDetails.Weight * float32(storedCartItemDetails.Quantity) * float32(fee.CostPerKg)
		newItemCost := storedCartItemDetails.Weight * float32(newQuantity) * float32(fee.CostPerKg)

		updateRequest.CartItemID = request.CartItemID
		updateRequest.Quantity = request.Quantity
		updateRequest.ItemTotalCost = float64(newItemCost)
		updateRequest.CartTotalCost = float64(newItemCost - oldItemCost)
	}

	return resp, &updateRequest, nil
}
