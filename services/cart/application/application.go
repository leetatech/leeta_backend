package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg/query"
	"time"

	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/mailer"
	"github.com/leetatech/leeta_backend/services/cart/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type CartAppHandler struct {
	idGenerator   pkg.IDGenerator
	tokenHandler  pkg.TokenHandler
	logger        *zap.Logger
	EmailClient   mailer.MailerClient
	allRepository pkg.Repositories
}

type CartApplication interface {
	DeleteCart(ctx context.Context, cartId string) (*pkg.DefaultResponse, error)
	DeleteCartItem(ctx context.Context, cartItemId string) (*pkg.DefaultResponse, error)
	AddToCart(ctx context.Context, request domain.CartItem) (*pkg.DefaultResponse, error)
	UpdateCartItemQuantity(ctx context.Context, request domain.UpdateCartItemQuantityRequest) (*pkg.DefaultResponse, error)
	ListCartItems(ctx context.Context, request *query.ResultSelector) (*domain.ListCartResponse, error)
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

func (c CartAppHandler) DeleteCart(ctx context.Context, cartId string) (*pkg.DefaultResponse, error) {
	err := c.allRepository.CartRepository.DeleteCart(ctx, cartId)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error deleting cart: %w", err))
	}
	return &pkg.DefaultResponse{Success: "success", Message: "Cart deleted successfully"}, nil
}

func (c CartAppHandler) AddToCart(ctx context.Context, request domain.CartItem) (*pkg.DefaultResponse, error) {
	claims, err := c.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting user claims %w", leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err))
	}

	product, err := c.allRepository.ProductRepository.GetProductByID(ctx, request.ProductID)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.InvalidProductIdError, fmt.Errorf("error getting product id %s: %w", request.ProductID, err))
	}

	switch product.ParentCategory {
	case models.LNGProductCategory, models.LPGProductCategory:
		if request.Weight == 0 {
			return nil, leetError.ErrorResponseBody(leetError.InvalidRequestError, errors.New("invalid cart item, cart weight cannot be zero"))
		}
	}

	fee, err := c.allRepository.FeesRepository.GetFeeByProductID(ctx, product.ID, models.FeesActive)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.InvalidProductIdError, fmt.Errorf("error getting fee %w", err))
	}

	cartItem := models.CartItem{
		ID:              c.idGenerator.Generate(),
		ProductID:       request.ProductID,
		ProductCategory: product.ParentCategory,
		VendorID:        product.VendorID,
		Weight:          request.Weight,
		Quantity:        request.Quantity,
	}

	cartItem.Cost, err = cartItem.CalculateCartItemFee(fee)
	if cartItem.Cost == 0 {
		return nil, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("unable to calculate cart fee %w", err))
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
				return nil, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error when adding item to cart store %w", addToCartErr))
			}
			return &pkg.DefaultResponse{Success: "success", Message: "Successfully added item to cart"}, nil
		default:
			return nil, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error getting cart item by customer id %w", err))
		}
	}

	cart.CartItems = append(cart.CartItems, cartItem)

	cart.Total, err = c.calculateCartItemTotalCost(ctx, cart.CartItems)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error calculating cart item total fee %w", err))
	}

	cart.StatusTs = time.Now().Unix()
	err = c.allRepository.CartRepository.AddToCartItem(ctx, cart.ID, cartItem, cart.Total, cart.StatusTs)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error adding item to cart %w", err))
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
				cartTotalFee, err := item.CalculateCartItemFee(&fee)
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

	cart, err := c.allRepository.CartRepository.GetCartByCartItemID(ctx, request.CartItemID)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error getting cart item with cart item id '%s': %w", request.CartItemID, err))
	}
	cartItem, index, err := c.retrieveCartItemFromUpdateRequest(ctx, request, cart)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error retrieving cartItem: %w", err))
	}
	adjustedCartItem, err := c.adjustCartItemAndCalculateCost(ctx, cartItem)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error adjust cart item and calculating cost: %w", err))
	}
	// update item in cart and calculate total cart new cost
	cart.CartItems[index] = adjustedCartItem
	totalCost := cart.CalculateCartTotalFee()
	cart.Total = totalCost

	err = c.allRepository.CartRepository.UpdateCart(ctx, *cart)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, fmt.Errorf("error updating cart with adjusted cart item: %w", err))
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Successfully updated cart item quantity"}, nil
}
func (c CartAppHandler) retrieveCartItemFromUpdateRequest(ctx context.Context, request domain.UpdateCartItemQuantityRequest, cart *models.Cart) (cartItem models.CartItem, index int, err error) {
	for i, item := range cart.CartItems {
		if item.ID == request.CartItemID {
			index = i
			cartItem = item
			cartItem.Quantity = request.Quantity
			break
		}
	}

	if cartItem.ProductID != "" {
		product, getProductErr := c.allRepository.ProductRepository.GetProductByID(ctx, cartItem.ProductID)
		if getProductErr != nil {
			return cartItem, index, fmt.Errorf("error retrieving product %w", getProductErr)
		}

		cartItem.ProductCategory = product.ParentCategory
	}

	return
}

func (c CartAppHandler) retrieveProductFee(ctx context.Context, productID string) (*models.Fee, error) {
	fee, err := c.allRepository.FeesRepository.GetFeeByProductID(ctx, productID, models.FeesActive)
	if err != nil {
		return nil, fmt.Errorf("error retrieving fee for product with id '%s': %w", productID, err)
	}

	return fee, nil
}

func (c CartAppHandler) adjustCartItemAndCalculateCost(ctx context.Context, item models.CartItem) (cartItem models.CartItem, err error) {
	fee, err := c.retrieveProductFee(ctx, item.ProductID)
	if err != nil {
		err = fmt.Errorf("error retriving cart item product fee %w", err)
		return
	}

	switch item.ProductCategory {
	case models.LNGProductCategory, models.LPGProductCategory:
		itemWeightCost := item.Weight * float32(fee.CostPerKg)
		itemCost := itemWeightCost * float32(item.Quantity)
		item.Cost = float64(itemCost)
	default:
		// TODO: refactor this when we know more product categories
		err = errors.New("invalid product category, when updating product quantity")
		return
	}

	return item, nil
}

func (c CartAppHandler) DeleteCartItem(ctx context.Context, itemId string) (*pkg.DefaultResponse, error) {
	_, err := c.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	cart, err := c.allRepository.CartRepository.GetCartByCartItemID(ctx, itemId)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, fmt.Errorf("error getting cart item with cart item id '%s': %w", itemId, err))
	}

	var itemCost float64
	for _, item := range cart.CartItems {
		if item.ID == itemId {
			itemCost = item.Cost
		}
	}

	err = c.allRepository.CartRepository.DeleteCartItem(ctx, itemId, itemCost)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, fmt.Errorf("error deleting cart item with cart item id '%s': %w", itemId, err))
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Successfully deleted item from cart"}, nil
}

func (c CartAppHandler) ListCartItems(ctx context.Context, request *query.ResultSelector) (*domain.ListCartResponse, error) {
	claims, err := c.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	cart, err := c.allRepository.CartRepository.ListCart(ctx, request, claims.UserID)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, fmt.Errorf("error listing cart items: %w", err))
	}

	return cart, nil
}
