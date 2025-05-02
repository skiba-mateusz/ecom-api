package repository

import (
	"context"
	"github.com/skiba-mateusz/ecom-api/internal/app/domain"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

type MockCategoryRepository struct {
	mock.Mock
}

func (r *MockProductRepository) GetById(ctx context.Context, id int64) (*domain.Product, error) {
	args := r.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Product), args.Error(1)
}

func (r *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	args := r.Called(ctx, product)
	return args.Error(0)
}

func (r *MockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	args := r.Called(ctx, product)
	return args.Error(0)
}

func (r *MockProductRepository) Delete(ctx context.Context, id int64) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

func (r *MockProductRepository) SlugExists(ctx context.Context, candidate string) (bool, error) {
	args := r.Called(ctx, candidate)
	return args.Bool(0), args.Error(1)
}

func (r *MockProductRepository) List(ctx context.Context, q domain.PaginatedProductsQuery) ([]domain.ProductSummary, domain.Meta, error) {
	args := r.Called(ctx, q)

	var products []domain.ProductSummary
	if args.Get(0) != nil {
		products = args.Get(0).([]domain.ProductSummary)
	}

	var meta domain.Meta
	if args.Get(1) != nil {
		meta = args.Get(1).(domain.Meta)
	}

	return products, meta, args.Error(2)
}

func (r *MockCategoryRepository) GetById(ctx context.Context, id int64) (*domain.Category, error) {
	args := r.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Category), args.Error(1)
}
