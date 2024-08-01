package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/paging"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/pkg/helpers"
	"github.com/leetatech/leeta_backend/pkg/idgenerator"
	"github.com/leetatech/leeta_backend/pkg/jwtmiddleware"
	mailer "github.com/leetatech/leeta_backend/pkg/notification/mailer/aws"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"time"

	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/services/cart/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartApplicationManager struct {
	idgenerator       idgenerator.Generator
	jwtManager        jwtmiddleware.Manager
	EmailClient       mailer.Client
	repositoryManager pkg.RepositoryManager
}

type Cart interface {
	Delete(ctx context.Context, id string) error
	DeleteItem(ctx context.Context, itemId string) error
	Add(ctx context.Context, item domain.CartItem) (models.Cart, error)
	UpdateItemQuantity(ctx context.Context, itemQuantity domain.UpdateCartItemQuantityRequest) (models.Cart, error)
	ListCart(ctx context.Context, request query.ResultSelector) (models.Cart, uint64, error)
	Checkout(ctx context.Context, request domain.CartCheckoutRequest) (*pkg.DefaultResponse, error)
}

func New(applicationContext pkg.ApplicationContext) Cart {
	return &CartApplicationManager{
		idgenerator:       idgenerator.New(),
		jwtManager:        applicationContext.JwtManager,
		EmailClient:       applicationContext.MailClient,
		repositoryManager: applicationContext.RepositoryManager,
	}
}

func (c *CartApplicationManager) Delete(ctx context.Context, cartId string) error {
	err := c.repositoryManager.CartRepository.DeleteCart(ctx, cartId)
	if err != nil {
		return errs.Body(errs.InternalError, fmt.Errorf("error deleting cart: %w", err))
	}
	return nil
}

func (c *CartApplicationManager) Add(ctx context.Context, request domain.CartItem) (cart models.Cart, err error) {
	claims, err := c.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return cart, errs.Body(errs.ErrorUnauthorized, fmt.Errorf("error getting user claims %w", err))
	}

	product, err := c.repositoryManager.ProductRepository.Product(ctx, request.ProductID)
	if err != nil {
		return cart, errs.Body(errs.InvalidProductIdError, fmt.Errorf("error getting product id %s: %w", request.ProductID, err))
	}

	switch product.ParentCategory {
	case models.LNGProductCategory, models.LPGProductCategory:
		if request.Weight == 0 {
			return cart, errs.Body(errs.InvalidRequestError, errors.New("invalid cart item, cart weight cannot be zero"))
		}
	}

	fee, err := c.repositoryManager.FeesRepository.ByProductID(ctx, product.ID, models.FeesActive)
	if err != nil {
		return cart, errs.Body(errs.FeesError, fmt.Errorf("error getting fee %w", err))
	}

	cartItem := models.CartItem{
		ID:              c.idgenerator.Generate(),
		ProductID:       request.ProductID,
		ProductCategory: product.ParentCategory,
		VendorID:        product.VendorID,
		Weight:          request.Weight,
		Quantity:        request.Quantity,
	}

	cartItem.Cost, err = cartItem.CalculateCartItemFee(fee)
	if cartItem.Cost == 0 || err != nil {
		return cart, errs.Body(errs.InternalError, fmt.Errorf("unable to calculate cart fee %w", err))
	}

	cart, err = c.repositoryManager.CartRepository.GetCartByCustomerID(ctx, claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			cart = models.Cart{
				ID:         c.idgenerator.Generate(),
				CustomerID: claims.UserID,
				CartItems:  []models.CartItem{cartItem},
				Total:      cartItem.Cost,
				Status:     models.CartActive,
				StatusTs:   time.Now().Unix(),
				Ts:         time.Now().Unix(),
			}

			addToCartErr := c.repositoryManager.CartRepository.AddToCart(ctx, cart)
			if addToCartErr != nil {
				return cart, errs.Body(errs.InternalError, fmt.Errorf("error when adding item to cart store %w", addToCartErr))
			}
			return cart, nil
		default:
			return cart, errs.Body(errs.InternalError, fmt.Errorf("error getting cart item by customer id %w", err))
		}
	}

	cart.CartItems = append(cart.CartItems, cartItem)

	cart.Total, err = c.calculateCartItemTotalCost(ctx, cart.CartItems)
	if err != nil {
		return cart, errs.Body(errs.InternalError, fmt.Errorf("error calculating cart item total fee %w", err))
	}

	cart.StatusTs = time.Now().Unix()
	err = c.repositoryManager.CartRepository.AddToCartItem(ctx, cart.ID, cartItem, cart.Total, cart.StatusTs)
	if err != nil {
		return cart, errs.Body(errs.InternalError, fmt.Errorf("error adding item to cart %w", err))
	}

	return cart, nil
}

func (c *CartApplicationManager) calculateCartItemTotalCost(ctx context.Context, items []models.CartItem) (float64, error) {
	var total float64

	fees, err := c.repositoryManager.FeesRepository.FeesByStatus(ctx, models.FeesActive)
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
func (c *CartApplicationManager) UpdateItemQuantity(ctx context.Context, request domain.UpdateCartItemQuantityRequest) (updatedCart models.Cart, err error) {
	_, err = c.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return updatedCart, errs.Body(errs.ErrorUnauthorized, err)
	}

	cart, err := c.repositoryManager.CartRepository.GetCartByCartItemID(ctx, request.CartItemID)
	if err != nil {
		return updatedCart, errs.Body(errs.InternalError, fmt.Errorf("error getting cart item with cart item id '%s': %w", request.CartItemID, err))
	}
	cartItem, index, err := c.retrieveCartItemFromUpdateRequest(ctx, request, cart)
	if err != nil {
		return updatedCart, errs.Body(errs.InternalError, fmt.Errorf("error retrieving cartItem: %w", err))
	}

	adjustedCartItem, err := c.adjustCartItemAndCalculateCost(ctx, cartItem)
	if err != nil {
		return updatedCart, errs.Body(errs.InternalError, fmt.Errorf("error adjust cart item and calculating cost: %w", err))
	}
	// update item in cart and calculate total cart new cost
	cart.CartItems[index] = adjustedCartItem
	totalCost := cart.CalculateCartTotalFee()
	cart.Total = totalCost

	err = c.repositoryManager.CartRepository.UpdateCart(ctx, cart)
	if err != nil {
		return updatedCart, errs.Body(errs.DatabaseError, fmt.Errorf("error updating cart with adjusted cart item: %w", err))
	}

	return cart, nil
}
func (c *CartApplicationManager) retrieveCartItemFromUpdateRequest(ctx context.Context, request domain.UpdateCartItemQuantityRequest, cart models.Cart) (cartItem models.CartItem, index int, err error) {
	for i, item := range cart.CartItems {
		if item.ID == request.CartItemID {
			index = i
			cartItem = item
			cartItem.Quantity = request.Quantity
			break
		}
	}

	if cartItem.ProductID != "" {
		product, getProductErr := c.repositoryManager.ProductRepository.Product(ctx, cartItem.ProductID)
		if getProductErr != nil {
			return cartItem, index, fmt.Errorf("error retrieving product %w", getProductErr)
		}

		cartItem.ProductCategory = product.ParentCategory
	}

	return
}

func (c *CartApplicationManager) retrieveProductFee(ctx context.Context, productID string) (*models.Fee, error) {
	fee, err := c.repositoryManager.FeesRepository.ByProductID(ctx, productID, models.FeesActive)
	if err != nil {
		return nil, fmt.Errorf("error retrieving fee for product with id '%s': %w", productID, err)
	}

	return fee, nil
}

func (c *CartApplicationManager) adjustCartItemAndCalculateCost(ctx context.Context, item models.CartItem) (cartItem models.CartItem, err error) {
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

func (c *CartApplicationManager) DeleteItem(ctx context.Context, itemId string) error {
	_, err := c.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return errs.Body(errs.ErrorUnauthorized, err)
	}

	cart, err := c.repositoryManager.CartRepository.GetCartByCartItemID(ctx, itemId)
	if err != nil {
		return errs.Body(errs.DatabaseError, fmt.Errorf("error getting cart item with cart item id '%s': %w", itemId, err))
	}

	var itemCost float64
	for _, item := range cart.CartItems {
		if item.ID == itemId {
			itemCost = item.Cost
		}
	}

	err = c.repositoryManager.CartRepository.DeleteCartItem(ctx, itemId, itemCost)
	if err != nil {
		return errs.Body(errs.DatabaseError, fmt.Errorf("error deleting cart item with cart item id '%s': %w", itemId, err))
	}

	return nil
}

func (c *CartApplicationManager) ListCart(ctx context.Context, request query.ResultSelector) (models.Cart, uint64, error) {
	claims, err := c.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return models.Cart{}, 0, errs.Body(errs.ErrorUnauthorized, err)
	}

	cart, totalResults, err := c.repositoryManager.CartRepository.ListCartItems(ctx, request, claims.UserID)
	if err != nil {
		return models.Cart{}, 0, errs.Body(errs.DatabaseError, fmt.Errorf("error listing cart items: %w", err))
	}

	return cart, totalResults, nil
}

func (c *CartApplicationManager) Checkout(ctx context.Context, request domain.CartCheckoutRequest) (*pkg.DefaultResponse, error) {
	claims, err := c.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return nil, errs.Body(errs.ErrorUnauthorized, err)
	}

	cart, err := c.repositoryManager.CartRepository.GetCartByCustomerID(ctx, claims.UserID)
	if err != nil {
		return nil, errs.Body(errs.DatabaseNoRecordError, err)
	}

	err = c.validateFees(ctx, request.DeliveryDetails.Address, request.DeliveryFee, request.ServiceFee)
	if err != nil {
		return nil, err
	}

	err = c.checkout(ctx, claims.UserID, request, cart)
	if err != nil {
		return nil, err
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Checkout successful"}, nil
}

func (c *CartApplicationManager) checkout(ctx context.Context, userID string, request domain.CartCheckoutRequest, cart models.Cart) (err error) {

	totalCost := cart.Total + request.DeliveryFee + request.ServiceFee

	if helpers.RoundToTwoDecimalPlaces(request.TotalFee) < helpers.RoundToTwoDecimalPlaces(totalCost) {
		return errs.Body(errs.AmountPaidError, errors.New("amount paid does not match total cost"))
	}

	orderStatus := []models.StatusHistory{
		{
			Status:   models.OrderPending,
			StatusTs: time.Now().Unix(),
		},
	}

	order := models.Order{
		ID:              c.idgenerator.Generate(),
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

	err = c.repositoryManager.OrderRepository.Create(ctx, order)
	if err != nil {
		return errs.Body(errs.InternalError, fmt.Errorf("error creating order when checking out of cart %w", err))
	}

	err = c.repositoryManager.CartRepository.ClearCart(ctx, cart.ID)
	if err != nil {
		return errs.Body(errs.InternalError, fmt.Errorf("error clearing cart %w", err))
	}

	return nil
}

func (c *CartApplicationManager) validateFees(ctx context.Context, address models.Address, deliveryFee, serviceFee float64) error {
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
	fees, _, err := c.repositoryManager.FeesRepository.Fees(ctx, getRequest)
	if err != nil {
		return err
	}
	// validate delivery fee
	for _, fee := range fees {
		switch fee.FeeType {
		case models.DeliveryFee:
			if fee.Cost.CostPerType != deliveryFee {
				return errs.Body(errs.InvalidDeliveryFeeError, errors.New("invalid delivery fee"))
			}

		case models.ServiceFee:
			if fee.Cost.CostPerType != serviceFee {
				return errs.Body(errs.InvalidServiceFeeError, errors.New("invalid service fee"))
			}
		}

	}
	return nil
}
