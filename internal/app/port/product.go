package port

import (
	"context"
	"github.com/skiba-mateusz/ecom-api/internal/app/domain"
)

type ProductRepository interface {
	GetById(ctx context.Context, id int64) (*domain.Product, error)
	Create(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, product *domain.Product) error
	SlugExists(ctx context.Context, candidate string) (bool, error)
	List(ctx context.Context, query domain.PaginatedProductsQuery) ([]domain.ProductSummary, domain.Meta, error)
}

type ProductService interface {
	GetById(ctx context.Context, id int64) (*domain.Product, error)
	Create(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, product *domain.Product) error
	List(ctx context.Context, query domain.PaginatedProductsQuery) ([]domain.ProductSummary, domain.Meta, error)
}
