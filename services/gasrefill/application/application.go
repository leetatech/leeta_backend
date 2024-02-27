package application

import (
	"context"
	"fmt"
	"github.com/leetatech/leeta_backend/services/gasrefill/domain"
	"github.com/leetatech/leeta_backend/services/library"
	"github.com/leetatech/leeta_backend/services/library/models"
	"go.uber.org/zap"
	"time"
)

type GasRefillHandler struct {
	idGenerator   library.IDGenerator
	logger        *zap.Logger
	allRepository library.Repositories
}

type GasRefillApplication interface {
	RequestRefill(request domain.GasRefillRequest) (*library.DefaultResponse, error)
}

func NewGasRefillApplication(request library.DefaultApplicationRequest) GasRefillApplication {
	return &GasRefillHandler{
		idGenerator:   library.NewIDGenerator(),
		logger:        request.Logger,
		allRepository: request.AllRepository,
	}
}

func (r *GasRefillHandler) RequestRefill(request domain.GasRefillRequest) (*library.DefaultResponse, error) {

	ctx := context.Background()
	if !request.Guest && request.CustomerID != "" {
		customer, err := r.allRepository.UserRepository.GetCustomerByID(request.CustomerID)
		if err != nil {
			return nil, err
		}

		if request.ShippingInfo.ForMe {
			request.ShippingInfo = r.forMeCheck(request.ShippingInfo, fmt.Sprintf("%s %s", customer.FirstName, customer.LastName), customer.Phone.Number, customer.Email.Address)
		}
	}

	if request.Guest && request.GuestBioData.Email != "" {
		if request.ShippingInfo.ForMe {
			request.ShippingInfo = r.forMeCheck(request.ShippingInfo, fmt.Sprintf("%s %s", request.GuestBioData.FirstName, request.GuestBioData.LastName), request.GuestBioData.Phone, request.GuestBioData.Email)
		}
	}

	product, err := r.allRepository.ProductRepository.GetProductByID(ctx, request.RefillDetails.ProductID)
	if err != nil {
		return nil, err
	}

	vendor, err := r.allRepository.UserRepository.GetVendorByID(product.VendorID)
	if err != nil {
		return nil, err
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
