package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg/helpers"
	"github.com/leetatech/leeta_backend/pkg/mailer/postmarkClient"
	"github.com/leetatech/leeta_backend/pkg/query"
	"github.com/leetatech/leeta_backend/pkg/query/filter"
	"github.com/leetatech/leeta_backend/pkg/query/paging"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"time"

	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/services/cart/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type CartAppHandler struct {
	idGenerator   pkg.IDGenerator
	tokenHandler  pkg.TokenHandler
	logger        *zap.Logger
	EmailClient   postmarkClient.MailerClient
	allRepository pkg.Repositories
}

type CartApplication interface {
	DeleteCart(ctx context.Context, cartId string) error
	DeleteCartItem(ctx context.Context, cartItemId string) error
	AddToCart(ctx context.Context, request domain.CartItem) (models.Cart, error)
	UpdateCartItemQuantity(ctx context.Context, request domain.UpdateCartItemQuantityRequest) (models.Cart, error)
	ListCart(ctx context.Context, request query.ResultSelector) (models.Cart, uint64, error)
	CartCheckout(ctx context.Context, request domain.CartCheckoutRequest) (*pkg.DefaultResponse, error)
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

func (c *CartAppHandler) DeleteCart(ctx context.Context, cartId string) error {
	err := c.allRepository.CartRepository.DeleteCart(ctx, cartId)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error deleting cart: %w", err))
	}
	return nil
}

func (c *CartAppHandler) AddToCart(ctx context.Context, request domain.CartItem) (cart models.Cart, err error) {
	claims, err := c.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return cart, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, fmt.Errorf("error getting user claims %w", err))
	}

	product, err := c.allRepository.ProductRepository.GetProductByID(ctx, request.ProductID)
	if err != nil {
		return cart, leetError.ErrorResponseBody(leetError.InvalidProductIdError, fmt.Errorf("error getting product id %s: %w", request.ProductID, err))
	}

	switch product.ParentCategory {
	case models.LNGProductCategory, models.LPGProductCategory:
		if request.Weight == 0 {
			return cart, leetError.ErrorResponseBody(leetError.InvalidRequestError, errors.New("invalid cart item, cart weight cannot be zero"))
		}
	}

	fee, err := c.allRepository.FeesRepository.GetFeeByProductID(ctx, product.ID, models.FeesActive)
	if err != nil {
		return cart, leetError.ErrorResponseBody(leetError.FeesError, fmt.Errorf("error getting fee %w", err))
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
	if cartItem.Cost == 0 || err != nil {
		return cart, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("unable to calculate cart fee %w", err))
	}

	cart, err = c.allRepository.CartRepository.GetCartByCustomerID(ctx, claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			cart = models.Cart{
				ID:         c.idGenerator.Generate(),
				CustomerID: claims.UserID,
				CartItems:  []models.CartItem{cartItem},
				Total:      cartItem.Cost,
				Status:     models.CartActive,
				StatusTs:   time.Now().Unix(),
				Ts:         time.Now().Unix(),
			}

			addToCartErr := c.allRepository.CartRepository.AddToCart(ctx, cart)
			if addToCartErr != nil {
				return cart, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error when adding item to cart store %w", addToCartErr))
			}
			return cart, nil
		default:
			return cart, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error getting cart item by customer id %w", err))
		}
	}

	cart.CartItems = append(cart.CartItems, cartItem)

	cart.Total, err = c.calculateCartItemTotalCost(ctx, cart.CartItems)
	if err != nil {
		return cart, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error calculating cart item total fee %w", err))
	}

	cart.StatusTs = time.Now().Unix()
	err = c.allRepository.CartRepository.AddToCartItem(ctx, cart.ID, cartItem, cart.Total, cart.StatusTs)
	if err != nil {
		return cart, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error adding item to cart %w", err))
	}

	return cart, nil
}

func (c *CartAppHandler) calculateCartItemTotalCost(ctx context.Context, items []models.CartItem) (float64, error) {
	var total float64

	fees, err := c.allRepository.FeesRepository.GetActiveFees(ctx, models.FeesActive)
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
func (c *CartAppHandler) UpdateCartItemQuantity(ctx context.Context, request domain.UpdateCartItemQuantityRequest) (updatedCart models.Cart, err error) {
	_, err = c.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return updatedCart, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	cart, err := c.allRepository.CartRepository.GetCartByCartItemID(ctx, request.CartItemID)
	if err != nil {
		return updatedCart, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error getting cart item with cart item id '%s': %w", request.CartItemID, err))
	}
	cartItem, index, err := c.retrieveCartItemFromUpdateRequest(ctx, request, cart)
	if err != nil {
		return updatedCart, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error retrieving cartItem: %w", err))
	}

	adjustedCartItem, err := c.adjustCartItemAndCalculateCost(ctx, cartItem)
	if err != nil {
		return updatedCart, leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error adjust cart item and calculating cost: %w", err))
	}
	// update item in cart and calculate total cart new cost
	cart.CartItems[index] = adjustedCartItem
	totalCost := cart.CalculateCartTotalFee()
	cart.Total = totalCost

	err = c.allRepository.CartRepository.UpdateCart(ctx, cart)
	if err != nil {
		return updatedCart, leetError.ErrorResponseBody(leetError.DatabaseError, fmt.Errorf("error updating cart with adjusted cart item: %w", err))
	}

	return cart, nil
}
func (c *CartAppHandler) retrieveCartItemFromUpdateRequest(ctx context.Context, request domain.UpdateCartItemQuantityRequest, cart models.Cart) (cartItem models.CartItem, index int, err error) {
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

func (c *CartAppHandler) retrieveProductFee(ctx context.Context, productID string) (*models.Fee, error) {
	fee, err := c.allRepository.FeesRepository.GetFeeByProductID(ctx, productID, models.FeesActive)
	if err != nil {
		return nil, fmt.Errorf("error retrieving fee for product with id '%s': %w", productID, err)
	}

	return fee, nil
}

func (c *CartAppHandler) adjustCartItemAndCalculateCost(ctx context.Context, item models.CartItem) (cartItem models.CartItem, err error) {
	fee, err := c.retrieveProductFee(ctx, item.ProductID)
	if err != nil {
		err = fmt.Errorf("error retriving cart item product fee %w", err)
		return
	}

	switch item.ProductCategory {
	case models.LNGProductCategory, models.LPGProductCategory:
		itemWeightCost := item.Weight * float32(fee.Cost.CostPerKG)
		itemCost := itemWeightCost * float32(item.Quantity)
		item.Cost = float64(itemCost)
	default:
		// TODO: refactor this when we know more product categories
		err = errors.New("invalid product category, when updating product quantity")
		return
	}

	return item, nil
}

func (c *CartAppHandler) DeleteCartItem(ctx context.Context, itemId string) error {
	_, err := c.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	cart, err := c.allRepository.CartRepository.GetCartByCartItemID(ctx, itemId)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, fmt.Errorf("error getting cart item with cart item id '%s': %w", itemId, err))
	}

	var itemCost float64
	for _, item := range cart.CartItems {
		if item.ID == itemId {
			itemCost = item.Cost
		}
	}

	err = c.allRepository.CartRepository.DeleteCartItem(ctx, itemId, itemCost)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, fmt.Errorf("error deleting cart item with cart item id '%s': %w", itemId, err))
	}

	return nil
}

func (c *CartAppHandler) ListCart(ctx context.Context, request query.ResultSelector) (models.Cart, uint64, error) {
	claims, err := c.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return models.Cart{}, 0, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	cart, totalResults, err := c.allRepository.CartRepository.ListCartItems(ctx, request, claims.UserID)
	if err != nil {
		return models.Cart{}, 0, leetError.ErrorResponseBody(leetError.DatabaseError, fmt.Errorf("error listing cart items: %w", err))
	}

	return cart, totalResults, nil
}

func (c *CartAppHandler) CartCheckout(ctx context.Context, request domain.CartCheckoutRequest) (*pkg.DefaultResponse, error) {
	claims, err := c.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	cart, err := c.allRepository.CartRepository.GetCartByCustomerID(ctx, claims.UserID)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseNoRecordError, err)
	}

	err = c.validateFees(ctx, request.DeliveryDetails.Address, request.DeliveryFee, request.ServiceFee)
	if err != nil {
		return nil, err
	}

	err = c.performCheckout(ctx, claims.UserID, request, cart)
	if err != nil {
		return nil, err
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Checkout successful"}, nil
}

func (c *CartAppHandler) performCheckout(ctx context.Context, userID string, request domain.CartCheckoutRequest, cart models.Cart) (err error) {

	totalCost := cart.Total + request.DeliveryFee + request.ServiceFee

	if helpers.RoundToTwoDecimalPlaces(request.TotalFee) < helpers.RoundToTwoDecimalPlaces(totalCost) {
		return leetError.ErrorResponseBody(leetError.AmountPaidError, errors.New("amount paid does not match total cost"))
	}

	orderStatus := []models.StatusHistory{
		{
			Status:   models.OrderPending,
			StatusTs: time.Now().Unix(),
		},
	}

	order := models.Order{
		ID:              c.idGenerator.Generate(),
		Orders:          cart.CartItems,
		CustomerID:      userID,
		DeliveryDetails: request.DeliveryDetails,
		PaymentMethod:   request.PaymentMethod,
		DeliveryFee:     request.DeliveryFee,
		ServiceFee:      request.ServiceFee,
		Total:           request.TotalFee,
		StatusHistory:   orderStatus,
		StatusTs:        time.Now().Unix(),
		Ts:              time.Now().Unix(),
	}

	err = c.allRepository.OrderRepository.CreateOrder(ctx, order)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error creating order when checking out of cart %w", err))
	}

	err = c.allRepository.CartRepository.ClearCart(ctx, cart.ID)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.InternalError, fmt.Errorf("error clearing cart %w", err))
	}

	return nil
}

func (c *CartAppHandler) validateFees(ctx context.Context, address models.Address, deliveryFee, serviceFee float64) error {
	// get delivery fee from database
	getRequest := query.ResultSelector{
		Filter: &filter.Request{
			Operator: "and",
			Fields: []filter.RequestField{
				{
					Name:  "lga",
					Value: models.LGA{LGA: address.LGA, State: strings.ToUpper(address.State)},
				},
				{
					Name:  "fee_type",
					Value: bson.M{"$in": []models.FeeType{models.DeliveryFee, models.ServiceFee}},
				},
				{
					Name:  "status",
					Value: models.FeesActive,
				},
			},
		},
		Paging: &paging.Request{},
	}
	fees, _, err := c.allRepository.FeesRepository.FetchFees(ctx, getRequest)
	if err != nil {
		return err
	}
	// validate delivery fee
	for _, fee := range fees {
		switch fee.FeeType {
		case models.DeliveryFee:
			if fee.Cost.CostPerType != deliveryFee {
				return leetError.ErrorResponseBody(leetError.InvalidDeliveryFeeError, errors.New("invalid delivery fee"))
			}

		case models.ServiceFee:
			if fee.Cost.CostPerType != serviceFee {
				return leetError.ErrorResponseBody(leetError.InvalidServiceFeeError, errors.New("invalid service fee"))
			}
		}

	}
	return nil
}
