package service

import (
	"context"
	"github.com/skiba-mateusz/ecom-api/internal/app/domain"
	"github.com/skiba-mateusz/ecom-api/internal/app/port"
	"github.com/skiba-mateusz/ecom-api/internal/app/util"
	"time"
)

type ProductService struct {
	productRepo  port.ProductRepository
	categoryRepo port.CategoryRepository
}

func NewProductService(productRepo port.ProductRepository, categoryRepo port.CategoryRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *ProductService) GetById(ctx context.Context, id int64) (*domain.Product, error) {
	product, err := s.productRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	category, err := s.categoryRepo.GetById(ctx, product.CategoryId)
	if err != nil {
		return nil, err
	}

	product.Category = category

	return product, nil
}

func (s *ProductService) Create(ctx context.Context, product *domain.Product) error {
	slug, err := util.GenerateUniqueSlug(ctx, product.Name, s.productRepo.SlugExists)
	if err != nil {
		return err
	}
	product.Slug = slug
	return s.productRepo.Create(ctx, product)
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	return s.productRepo.Delete(ctx, id)
}

func (s *ProductService) Update(ctx context.Context, product *domain.Product) error {
	existingProduct, err := s.productRepo.GetById(ctx, product.Id)
	if err != nil {
		return err
	}

	if existingProduct.Name != product.Name {
		slug, err := util.GenerateUniqueSlug(ctx, product.Name, s.productRepo.SlugExists)
		if err != nil {
			return err
		}
		product.Slug = slug
	} else {
		product.Slug = existingProduct.Slug
	}

	product.UpdatedAt = time.Now()

	return s.productRepo.Update(ctx, product)
}

func (s *ProductService) List(ctx context.Context, query domain.PaginatedProductsQuery) ([]domain.ProductSummary, domain.Meta, error) {
	return s.productRepo.List(ctx, query)
}
