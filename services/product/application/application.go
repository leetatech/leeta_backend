package application

import (
	"context"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/mailer"
	"github.com/leetatech/leeta_backend/pkg/query"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/product/domain"
	"go.uber.org/zap"
	"time"
)

type productAppHandler struct {
	tokenHandler  pkg.TokenHandler
	encryptor     pkg.EncryptorManager
	idGenerator   pkg.IDGenerator
	otpGenerator  pkg.OtpGenerator
	logger        *zap.Logger
	EmailClient   mailer.MailerClient
	allRepository pkg.Repositories
}

type ProductApplication interface {
	CreateProduct(ctx context.Context, request domain.ProductRequest) (*pkg.DefaultResponse, error)
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	GetAllVendorProducts(ctx context.Context, request domain.GetVendorProductsRequest) (*domain.GetVendorProductsResponse, error)
	ListProducts(ctx context.Context, request *query.ResultSelector) (*query.ResponseListWithMetadata[models.Product], error)
	CreateGasProduct(ctx context.Context, request domain.GasProductRequest) (*pkg.DefaultResponse, error)
}

// *query.ResponseListWithMetadata[CartResponseData]

func NewProductApplication(request pkg.DefaultApplicationRequest) ProductApplication {
	return &productAppHandler{
		tokenHandler:  request.TokenHandler,
		encryptor:     pkg.NewEncryptor(),
		idGenerator:   pkg.NewIDGenerator(),
		otpGenerator:  pkg.NewOTPGenerator(),
		logger:        request.Logger,
		EmailClient:   request.EmailClient,
		allRepository: request.AllRepository,
	}
}

func (p productAppHandler) CreateProduct(ctx context.Context, request domain.ProductRequest) (*pkg.DefaultResponse, error) {
	claims, err := p.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	if claims.Role != models.VendorCategory && claims.Role != models.AdminCategory {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	switch claims.Role {
	case models.AdminCategory:
		_, err = p.allRepository.AuthRepository.GetAdminByEmail(ctx, claims.Email)
		if err != nil {
			return nil, err
		}

	case models.VendorCategory:
		_, err = p.allRepository.UserRepository.GetVendorByID(request.VendorID)
		if err != nil {
			return nil, err
		}

	}
	finalPrice := request.OriginalPriceAndVat
	if request.Discount {
		finalPrice = (request.OriginalPrice - request.DiscountPrice) + request.Vat
	}

	product := models.Product{
		ID:                  p.idGenerator.Generate(),
		VendorID:            request.VendorID,
		SubCategory:         request.SubCategory,
		Images:              request.Images,
		Name:                request.Name,
		Weight:              request.Weight,
		Description:         request.Description,
		OriginalPrice:       request.OriginalPrice,
		Vat:                 request.Vat,
		OriginalPriceAndVat: request.OriginalPriceAndVat,
		Discount:            true,
		DiscountPrice:       request.DiscountPrice,
		FinalPrice:          finalPrice,
		Status:              request.Status,
		StatusTs:            time.Now().Unix(),
		Ts:                  time.Now().Unix(),
	}

	err = p.allRepository.ProductRepository.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Product successfully created"}, nil
}

func (p productAppHandler) CreateGasProduct(ctx context.Context, request domain.GasProductRequest) (*pkg.DefaultResponse, error) {
	_, err := p.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	productCategory, err := models.SetProductCategory(request.ProductCategory)
	if err != nil {
		return nil, err
	}

	product := models.Product{
		ID:             p.idGenerator.Generate(),
		Name:           request.Name,
		ParentCategory: productCategory,
		Description:    request.Description,
		Status:         models.InStock,
		StatusTs:       time.Now().Unix(),
		Ts:             time.Now().Unix(),
	}

	err = p.allRepository.ProductRepository.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Gas Product successfully created"}, nil
}

func (p productAppHandler) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	_, err := p.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}
	product, err := p.allRepository.ProductRepository.GetProductByID(ctx, id)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return product, nil
}

func (p productAppHandler) GetAllVendorProducts(ctx context.Context, request domain.GetVendorProductsRequest) (*domain.GetVendorProductsResponse, error) {
	_, err := p.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	products, err := p.allRepository.ProductRepository.GetAllVendorProducts(ctx, request)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p productAppHandler) ListProducts(ctx context.Context, request *query.ResultSelector) (*query.ResponseListWithMetadata[models.Product], error) {
	_, err := p.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, leetError.ErrorResponseBody(leetError.ErrorUnauthorized, err)
	}

	products, err := p.allRepository.ProductRepository.ListProducts(ctx, request)
	if err != nil {
		return nil, err
	}

	return products, nil
}
