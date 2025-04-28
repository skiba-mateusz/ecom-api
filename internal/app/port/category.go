package port

import (
	"context"
	"github.com/skiba-mateusz/ecom-api/internal/app/domain"
)

type CategoryRepository interface {
	GetById(ctx context.Context, id int64) (*domain.Category, error)
}
