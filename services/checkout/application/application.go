package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/helpers"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/mailer"
	"github.com/leetatech/leeta_backend/pkg/query"
	"github.com/leetatech/leeta_backend/pkg/query/filter"
	"github.com/leetatech/leeta_backend/pkg/query/paging"
	"github.com/leetatech/leeta_backend/services/checkout/domain"
	"github.com/leetatech/leeta_backend/services/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"strings"
	"time"
)

type CheckoutHandler struct {
	idGenerator   pkg.IDGenerator
	tokenHandler  pkg.TokenHandler
	logger        *zap.Logger
	EmailClient   mailer.MailerClient
	allRepository pkg.Repositories
}

type CheckoutApplication interface {
	Checkout(ctx context.Context, request domain.CheckoutRequest) (*pkg.DefaultResponse, error)
	UpdateCheckout(ctx context.Context, request domain.UpdateCheckoutRequest) (*pkg.DefaultResponse, error)
}

func NewCheckoutApplication(request pkg.DefaultApplicationRequest) CheckoutApplication {
	return &CheckoutHandler{
		idGenerator:   pkg.NewIDGenerator(),
		logger:        request.Logger,
		tokenHandler:  request.TokenHandler,
		EmailClient:   request.EmailClient,
		allRepository: request.AllRepository,
	}
}

func (c *CheckoutHandler) Checkout(ctx context.Context, request domain.CheckoutRequest) (*pkg.DefaultResponse, error) {

	claims, err := c.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	cart, err := c.allRepository.CartRepository.GetCartByCustomerID(ctx, claims.UserID)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseNoRecordError, err)
	}

	defaultDeliveryDetails, err := c.setDefaultDeliveryDetails(ctx, claims, request.DeliveryDetails)
	if err != nil {
		return nil, err
	}

	var defaultDetails models.ShippingInfo
	if defaultDeliveryDetails != defaultDetails {
		request.DeliveryDetails = defaultDeliveryDetails
	}

	err = c.validateFees(ctx, request.DeliveryDetails.RecipientAddress, request.DeliveryFee, request.ServiceFee)
	if err != nil {
		return nil, err
	}

	err = c.performCheckout(ctx, claims.UserID, request, cart)
	if err != nil {
		return nil, err
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Checkout successful"}, nil
}

func (c *CheckoutHandler) performCheckout(ctx context.Context, userID string, request domain.CheckoutRequest, cart models.Cart) error {

	totalCost := cart.Total + request.DeliveryFee + request.ServiceFee

	if helpers.RoundToTwoDecimalPlaces(request.AmountPaid) < helpers.RoundToTwoDecimalPlaces(totalCost) {
		return leetError.ErrorResponseBody(leetError.AmountPaidError, errors.New("amount paid does not match total cost"))
	}

	checkout := models.Checkout{
		ID:         c.idGenerator.Generate(),
		CustomerID: userID,
		CheckoutDetails: models.CheckoutDetails{
			CartItems: cart.CartItems,
			Status:    models.CheckoutPending,
			StatusTs:  time.Now().Unix(),
			Ts:        time.Now().Unix(),
		},
		ShippingInfo: request.DeliveryDetails,
		AmountPaid:   request.AmountPaid,
		DeliveryFee:  request.DeliveryFee,
		ServiceFee:   request.ServiceFee,
		TotalCost:    helpers.RoundToTwoDecimalPlaces(totalCost),
		Status:       models.CheckoutPending,
		StatusTs:     time.Now().Unix(),
		Ts:           time.Now().Unix(),
	}

	err := c.allRepository.CheckoutRepository.RequestCheckout(ctx, checkout)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	err = c.allRepository.CartRepository.ClearCart(ctx, cart.ID)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}

func (c *CheckoutHandler) setDefaultDeliveryDetails(ctx context.Context, claims *pkg.UserClaims, deliveryDetails models.ShippingInfo) (models.ShippingInfo, error) {
	if deliveryDetails.ForMe {
		switch claims.Role {
		case models.GuestCatergory:

			guest, err := c.allRepository.AuthRepository.GetGuestRecord(ctx, claims.DeviceID)
			if err != nil {
				return models.ShippingInfo{}, leetError.ErrorResponseBody(leetError.DatabaseError, err)
			}

			deliveryDetails.RecipientName = guest.Name
			deliveryDetails.RecipientPhone = guest.Number
			deliveryDetails.RecipientEmail = guest.Email
			deliveryDetails.RecipientAddress = guest.Address
		case models.BuyerCategory:
			customer, err := c.allRepository.UserRepository.GetCustomerByID(claims.UserID)
			if err != nil {
				return models.ShippingInfo{}, leetError.ErrorResponseBody(leetError.DatabaseError, err)
			}
			deliveryDetails.RecipientName = fmt.Sprintf("%s %s", customer.FirstName, customer.LastName)
			deliveryDetails.RecipientPhone = customer.Phone.Number
			deliveryDetails.RecipientEmail = customer.Email.Address
			deliveryDetails.RecipientAddress = customer.Address
		}
	} else {
		if deliveryDetails.RecipientName == "" || deliveryDetails.RecipientPhone == "" || deliveryDetails.RecipientEmail == "" || deliveryDetails.RecipientAddress == (models.Address{}) {
			return models.ShippingInfo{}, leetError.ErrorResponseBody(leetError.InvalidRequestError, errors.New("invalid recipient details"))
		}

	}

	return deliveryDetails, nil
}

func (c *CheckoutHandler) validateFees(ctx context.Context, address models.Address, deliveryFee, serviceFee float64) error {
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

func (c *CheckoutHandler) UpdateCheckout(ctx context.Context, request domain.UpdateCheckoutRequest) (*pkg.DefaultResponse, error) {
	claims, err := c.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	if claims.Role == models.VendorCategory || claims.Role == models.AdminCategory {

		err = c.allRepository.CheckoutRepository.UpdateCheckoutStatus(ctx, request.ID, request.Status)
		if err != nil {
			return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
		}

		return &pkg.DefaultResponse{Success: "success", Message: "Checkout status updated successfully"}, nil
	}

	status, err := models.SetCheckoutStatus(request.Status)
	if err != nil {
		return nil, err
	}

	if status != models.CheckoutCancelled {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, errors.New("you cannot update this status"))
	}

	err = c.allRepository.CheckoutRepository.UpdateCheckoutStatus(ctx, request.ID, status)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Checkout status updated successfully"}, nil
}
