package application

import (
	"context"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/messaging/mailer/postmarkClient"
	"github.com/leetatech/leeta_backend/pkg/encrypto"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/pkg/idgenerator"
	"github.com/leetatech/leeta_backend/pkg/jwtmiddleware"
	"github.com/leetatech/leeta_backend/pkg/mailer/aws"
	"github.com/leetatech/leeta_backend/pkg/otp"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/product/domain"
	"time"
)

type productAppHandler struct {
	jwtManager    jwtmiddleware.Manager
	encryptor     encrypto.Manager
	idGenerator   idgenerator.Generator
	otpGenerator  otp.Generator
	EmailClient   aws.MailClient
	allRepository pkg.RepositoryManager
}

type Product interface {
	Create(ctx context.Context, request domain.ProductRequest) (*pkg.DefaultResponse, error)
	ProductByID(ctx context.Context, id string) (models.Product, error)
	VendorProducts(ctx context.Context, request domain.GetVendorProductsRequest) ([]models.Product, error)
	Products(ctx context.Context, request query.ResultSelector) (products []models.Product, totalResults uint64, err error)
	CreateGas(ctx context.Context, request domain.GasProductRequest) (*pkg.DefaultResponse, error)
}

func New(request pkg.ApplicationContext) Product {
	return &productAppHandler{
		jwtManager:    request.JwtManager,
		encryptor:     encrypto.New(),
		idGenerator:   idgenerator.New(),
		otpGenerator:  otp.New(),
		EmailClient:   request.MailClient,
		allRepository: request.RepositoryManager,
	}
}

func (p *productAppHandler) Create(ctx context.Context, request domain.ProductRequest) (*pkg.DefaultResponse, error) {
	claims, err := p.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return nil, errs.Body(errs.ErrorUnauthorized, err)
	}

	if claims.Role != models.VendorCategory && claims.Role != models.AdminCategory {
		return nil, errs.Body(errs.ErrorUnauthorized, err)
	}

	switch claims.Role {
	case models.AdminCategory:
		_, err = p.allRepository.AuthRepository.AdminByEmail(ctx, claims.Email)
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

	err = p.allRepository.ProductRepository.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Product successfully created"}, nil
}

func (p *productAppHandler) CreateGas(ctx context.Context, request domain.GasProductRequest) (*pkg.DefaultResponse, error) {
	_, err := p.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return nil, errs.Body(errs.ErrorUnauthorized, err)
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

	err = p.allRepository.ProductRepository.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return &pkg.DefaultResponse{Success: "success", Message: "Gas Product successfully created"}, nil
}

func (p *productAppHandler) ProductByID(ctx context.Context, id string) (product models.Product, err error) {
	_, err = p.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return product, errs.Body(errs.ErrorUnauthorized, err)
	}
	product, err = p.allRepository.ProductRepository.Product(ctx, id)
	if err != nil {
		return product, errs.Body(errs.DatabaseError, err)
	}

	return product, nil
}

func (p *productAppHandler) VendorProducts(ctx context.Context, request domain.GetVendorProductsRequest) ([]models.Product, error) {
	_, err := p.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return nil, errs.Body(errs.ErrorUnauthorized, err)
	}

	products, err := p.allRepository.ProductRepository.VendorProducts(ctx, request)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p *productAppHandler) Products(ctx context.Context, request query.ResultSelector) (products []models.Product, totalResults uint64, err error) {
	_, err = p.jwtManager.ExtractUserClaims(ctx)
	if err != nil {
		return nil, 0, errs.Body(errs.ErrorUnauthorized, err)
	}

	products, totalResults, err = p.allRepository.ProductRepository.ListProducts(ctx, request)
	if err != nil {
		return nil, 0, err
	}

	return products, totalResults, nil
}
