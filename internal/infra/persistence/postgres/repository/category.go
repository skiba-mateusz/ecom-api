package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/skiba-mateusz/ecom-api/internal/app/domain"
	"github.com/skiba-mateusz/ecom-api/internal/infra/persistence/postgres"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db}
}

func (r *CategoryRepository) GetById(ctx context.Context, id int64) (*domain.Category, error) {
	query := `
		WITH RECURSIVE tree AS (
		    SELECT 
		        id, name, slug, description, parent_id, image_url
			FROM categories WHERE id = $1 AND is_active = true
			UNION ALL
			SELECT 
			    c.id, c.name, c.slug, c.description, c.parent_id, c.image_url
			FROM categories c
			JOIN tree t ON c.id = t.parent_id
			WHERE c.is_active = true
		    )
		SELECT 
		    id, name, slug, description, parent_id, image_url
		FROM tree;
	`

	ctx, cancel := context.WithTimeout(ctx, postgres.QueryTimeoutDuration)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrNotFound
		default:
			return nil, err
		}
	}
	defer rows.Close()

	var categories []*domain.Category
	categoryMap := map[int64]*domain.Category{}
	for rows.Next() {
		var category domain.Category
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.ParentId,
			&category.ImageUrl,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
		categoryMap[category.Id] = &category
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return nil, domain.ErrNotFound
	}

	for _, category := range categories {
		if category.ParentId == nil {
			continue
		}

		if parent, exists := categoryMap[*category.ParentId]; exists {
			category.Parent = parent
		}
	}

	requestedCategory, exists := categoryMap[id]
	if !exists {
		return nil, domain.ErrNotFound
	}

	return requestedCategory, nil
}
