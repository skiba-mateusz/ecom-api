package handler

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/skiba-mateusz/ecom-api/internal/app/domain"
	"github.com/skiba-mateusz/ecom-api/internal/app/port"
	"github.com/skiba-mateusz/ecom-api/internal/infra/config"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type productIdKey string

const productIdCtx productIdKey = "productId"

type ProductHandler struct {
	config         *config.Config
	logger         *zap.SugaredLogger
	productService port.ProductService
}

func NewProductHandler(config *config.Config, logger *zap.SugaredLogger, productService port.ProductService) *ProductHandler {
	return &ProductHandler{
		config:         config,
		logger:         logger,
		productService: productService,
	}
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := getProductIdFromCtx(r.Context())

	product, err := h.productService.GetById(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			notFoundResponse(w, r, err, h.logger)
		default:
			internalServerError(w, r, err, h.logger)
		}
		return
	}

	if err = jsonResponse(w, http.StatusOK, product); err != nil {
		internalServerError(w, r, err, h.logger)
	}
}

type createProductRequest struct {
	Name        string   `json:"name" validate:"required,min=6,max=255"`
	Description *string  `json:"description" validate:"omitempty,min=32,max=1000"`
	Stock       int64    `json:"stock" validate:"required,min=0"`
	Price       float64  `json:"price" validate:"required,min=1"`
	SalePrice   *float64 `json:"sale_price" validate:"omitempty,min=1"`
	CategoryID  int64    `json:"category_id" validate:"required,min=1"`
	BrandID     int64    `json:"brand_id" validate:"required,min=1"`
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req createProductRequest
	if err := readJSON(w, r, &req); err != nil {
		badRequestResponse(w, r, err, h.logger)
		return
	}

	if err := validate.Struct(&req); err != nil {
		badRequestResponse(w, r, err, h.logger)
		return
	}

	product := &domain.Product{
		BaseProduct: domain.BaseProduct{
			Name:       req.Name,
			Stock:      req.Stock,
			Price:      req.Price,
			SalePrice:  req.SalePrice,
			CategoryId: req.CategoryID,
			BrandId:    req.BrandID,
		},
		Description: req.Description,
	}

	if err := h.productService.Create(r.Context(), product); err != nil {
		internalServerError(w, r, err, h.logger)
		return
	}

	if err := jsonResponse(w, http.StatusCreated, product); err != nil {
		internalServerError(w, r, err, h.logger)
	}
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := getProductIdFromCtx(r.Context())

	if err := h.productService.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			notFoundResponse(w, r, err, h.logger)
		default:
			internalServerError(w, r, err, h.logger)
		}
		return
	}

	if err := jsonResponse(w, http.StatusNoContent, nil); err != nil {
		internalServerError(w, r, err, h.logger)
	}
}

type updateProductRequest struct {
	Name        string   `json:"name" validate:"min=6,max=255"`
	Description *string  `json:"description" validate:"omitempty,required,min=32,max=1000"`
	Stock       int64    `json:"stock" validate:"required,min=0"`
	Price       float64  `json:"price" validate:"required,min=1"`
	SalePrice   *float64 `json:"sale_price" validate:"omitempty,min=1"`
	CategoryID  int64    `json:"category_id" validate:"required,min=1"`
	BrandID     int64    `json:"brand_id" validate:"required,min=1"`
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := getProductIdFromCtx(r.Context())

	var req updateProductRequest
	if err := readJSON(w, r, &req); err != nil {
		badRequestResponse(w, r, err, h.logger)
		return
	}

	if err := validate.Struct(&req); err != nil {
		badRequestResponse(w, r, err, h.logger)
		return
	}

	product := &domain.Product{
		BaseProduct: domain.BaseProduct{
			Id:         id,
			Name:       req.Name,
			Stock:      req.Stock,
			Price:      req.Price,
			SalePrice:  req.SalePrice,
			CategoryId: req.CategoryID,
			BrandId:    req.BrandID,
		},
		Description: req.Description,
	}

	if err := h.productService.Update(r.Context(), product); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			notFoundResponse(w, r, err, h.logger)
		default:
			internalServerError(w, r, err, h.logger)
		}
		return
	}

	if err := jsonResponse(w, http.StatusOK, product); err != nil {
		internalServerError(w, r, err, h.logger)
	}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	query := domain.PaginatedProductsQuery{
		Offset:        0,
		Limit:         20,
		Search:        "",
		SortField:     "name",
		SortDirection: "desc",
	}

	query, err := query.Parse(r)
	if err != nil {
		badRequestResponse(w, r, err, h.logger)
		return
	}

	if err = validate.Struct(query); err != nil {
		badRequestResponse(w, r, err, h.logger)
		return
	}

	products, meta, err := h.productService.List(r.Context(), query)
	if err != nil {
		internalServerError(w, r, err, h.logger)
		return
	}

	productsWithMeta := struct {
		Meta     domain.Meta             `json:"meta"`
		Products []domain.ProductSummary `json:"products"`
	}{
		Meta:     meta,
		Products: products,
	}

	if err = jsonResponse(w, http.StatusOK, productsWithMeta); err != nil {
		internalServerError(w, r, err, h.logger)
	}
}

func (h *ProductHandler) ProductIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			internalServerError(w, r, err, h.logger)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, productIdCtx, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getProductIdFromCtx(ctx context.Context) int64 {
	val := ctx.Value(productIdCtx)
	if val == nil {
		return 0
	}
	return val.(int64)
}
