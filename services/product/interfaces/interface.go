package interfaces

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	_ "github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/leetatech/leeta_backend/pkg"
	"github.com/leetatech/leeta_backend/pkg/helpers"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/product/application"
	"github.com/leetatech/leeta_backend/services/product/domain"
	"github.com/leetatech/leeta_backend/services/web"
	"github.com/samber/lo"
	"net/http"
)

type ProductHttpHandler struct {
	ProductApplication application.ProductApplication
}

func NewProductHTTPHandler(productApplication application.ProductApplication) *ProductHttpHandler {
	return &ProductHttpHandler{
		ProductApplication: productApplication,
	}

}

// CreateProductHandler godoc
// @Summary Create Product
// @Description The endpoint takes the product request and creates a new product
// @Tags Product
// @Accept multipart/form-data
// @Produce json
// @Param vendor_id formData string true "Vendor ID"
// @Param parent_category formData string false "Product parent category"
// @Param sub_category formData string true "Product subcategory"
// @Param name formData string true "Product name"
// @Param weight formData string true "Product weight"
// @Param description formData string true "Product description"
// @Param original_price formData string true "Product Price"
// @Param vat formData string true "Product vat"
// @Param original_price_and_vat formData string true "Product vat with original price"
// @Param discount formData string true "product discount availability"
// @Param discount_price formData string true "discount price"
// @Param status formData string true "product status"
// @Param images formData file true "Images of the product" format(multi)
// @Security BearerToken
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /product/create [post]
// @deprecated
func (handler *ProductHttpHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	request, err := checkFormFileAndAddProducts(r)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}

	resp, err := handler.ProductApplication.CreateProduct(r.Context(), *request)
	if err != nil {
		helpers.CheckErrorType(err, w)
		return
	}
	pkg.EncodeResult(w, resp, http.StatusOK)
}

// GetProductByIDHandler godoc
// @Summary Get Vendor Product By id
// @Description The endpoint takes the product id and then returns the requested product
// @Tags Product
// @Accept json
// @produce json
// @Param			product_id	path		string	true	"product id"
// @Security BearerToken
// @success 200 {object} models.Product
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /product/id/{product_id} [get]
// @deprecated
func (handler *ProductHttpHandler) GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	var (
		product models.Product
		err     error
	)
	productID := chi.URLParam(r, "product_id")

	product, err = handler.ProductApplication.GetProductByID(r.Context(), productID)
	if err != nil {
		helpers.CheckErrorType(err, w)
		return
	}

	pkg.EncodeResult(w, product, http.StatusOK)
}

// GetAllVendorProductsHandler godoc
// @Summary Get All Vendor Products By Status
// @Description The endpoint takes the vendor ID, product status, pages and limit and then returns the requested products
// @Tags Product
// @Accept json
// @produce json
// @param domain.GetVendorProductsRequest body domain.GetVendorProductsRequest true "get all vendor products request body"
// @Security BearerToken
// @success 200 {object} domain.GetVendorProductsResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /product/ [get]
// @deprecated
func (handler *ProductHttpHandler) GetAllVendorProductsHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.GetVendorProductsRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, err, http.StatusBadRequest)
		return
	}
	products, err := handler.ProductApplication.GetAllVendorProducts(r.Context(), request)
	if err != nil {
		helpers.CheckErrorType(err, w)
		return
	}
	pkg.EncodeResult(w, products, http.StatusOK)
}

// CreateGasProductHandler godoc
// @Summary Create Gas Product
// @Description The endpoint takes the gas product request and creates a new gas product
// @Tags Product
// @Accept json
// @Produce json
// @param domain.GasProductRequest body domain.GasProductRequest true "create gas product request body"
// @Security BearerToken
// @Success 200 {object} pkg.DefaultResponse
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /product/ [post]
func (handler *ProductHttpHandler) CreateGasProductHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.GasProductRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkg.EncodeResult(w, leetError.ErrorResponseBody(leetError.UnmarshalError, err), http.StatusBadRequest)
		return
	}

	result, err := handler.ProductApplication.CreateGasProduct(r.Context(), request)
	if err != nil {
		helpers.CheckErrorType(err, w)
		return
	}
	pkg.EncodeResult(w, result, http.StatusOK)
}

// ListProductsHandler godoc
// @Summary List Products
// @Description The endpoint takes in the limit, page and product status and returns the requested products
// @Tags Product
// @Accept json
// @Produce json
// @Param query.ResultSelector body query.ResultSelector true "list products request body"
// @Security BearerToken
// @Success 200 {object} query.ResponseListWithMetadata[models.Product]
// @Failure 401 {object} pkg.DefaultErrorResponse
// @Failure 400 {object} pkg.DefaultErrorResponse
// @Router /product/ [put]
func (handler *ProductHttpHandler) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	resultSelector, err := web.PrepareResultSelector(r, listProductOptions, allowedSortFields, web.ResultSelectorDefaults(defaultSortingRequest))
	if err != nil {
		pkg.EncodeErrorResult(w, http.StatusBadRequest, leetError.ErrorResponseBody(leetError.InvalidRequestError, err))
		return
	}

	products, totalResults, err := handler.ProductApplication.ListProducts(r.Context(), resultSelector)
	if err != nil {
		helpers.CheckErrorType(err, w)
		return
	}

	response := query.ResponseListWithMetadata[models.Product]{
		Metadata: query.NewMetadata(resultSelector, totalResults),
		Data:     products,
	}

	pkg.EncodeResult(w, response, http.StatusOK)
}

// ListProductOptions list product filter options
// @Summary Get Product filter options
// @Description Retrieve products filter options
// @Tags Product
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header" example(Bearer lnsjkfbnkjkdjnfjk)
// @Success 200 {object} []filter.RequestOption
// @Router /product/options [get]
func (handler *ProductHttpHandler) ListProductOptions(w http.ResponseWriter, r *http.Request) {
	requestOptions := lo.Map(listProductOptions, ToFilterOption)
	pkg.EncodeResult(w, requestOptions, http.StatusOK)
}
