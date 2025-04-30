package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/skiba-mateusz/ecom-api/internal/app/domain"
	"github.com/skiba-mateusz/ecom-api/internal/infra/persistence/postgres"
	"math"
	"strconv"
	"strings"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) GetById(ctx context.Context, id int64) (*domain.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.slug, p.description, p.price, p.sale_price, p.stock, p.category_id, p.brand_id, p.created_at, p.updated_at,
			b.id, b.name, b.slug, b.description, b.logo_url
		FROM products p
		LEFT JOIN brands b on p.brand_id = b.id
		WHERE p.id = $1 AND p.is_active = true;
	`

	ctx, cancel := context.WithTimeout(ctx, postgres.QueryTimeoutDuration)
	defer cancel()

	var product domain.Product
	product.Category = &domain.Category{}
	product.Brand = &domain.Brand{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.Id,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Price,
		&product.SalePrice,
		&product.Stock,
		&product.CategoryId,
		&product.BrandId,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.Brand.Id,
		&product.Brand.Name,
		&product.Brand.Slug,
		&product.Brand.Description,
		&product.Brand.LogoUrl,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrNotFound
		default:
			return nil, err
		}
	}

	return &product, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO 
		    products (name, slug, description, price, sale_price, stock, category_id, brand_id)
		VALUES 
		    ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			id, created_at, updated_at;
	`

	ctx, cancel := context.WithTimeout(ctx, postgres.QueryTimeoutDuration)
	defer cancel()

	err := r.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Slug,
		product.Description,
		product.Price,
		product.SalePrice,
		product.Stock,
		product.CategoryId,
		product.BrandId,
	).
		Scan(
			&product.Id,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE products SET is_active = false WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, postgres.QueryTimeoutDuration)
	defer cancel()

	res, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE 
		    products 
		SET
		    name = $1,
		    slug = $2,
		    description = $3,
		    price = $4,
		    sale_price = $5,
		    stock = $6,
		    category_id = $7,
		    brand_id = $8,
			updated_at = $9
		where 
		    id = $10
	`

	ctx, cancel := context.WithTimeout(ctx, postgres.QueryTimeoutDuration)
	defer cancel()

	res, err := r.db.ExecContext(
		ctx,
		query,
		product.Name,
		product.Slug,
		product.Description,
		product.Price,
		product.SalePrice,
		product.Stock,
		product.CategoryId,
		product.BrandId,
		product.UpdatedAt,
		product.Id,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *ProductRepository) SlugExists(ctx context.Context, candidate string) (bool, error) {
	query := `
		SELECT EXISTS (SELECT 1 FROM products WHERE slug = $1);
	`

	ctx, cancel := context.WithTimeout(ctx, postgres.QueryTimeoutDuration)
	defer cancel()

	var exists bool
	err := r.db.QueryRowContext(ctx, query, candidate).Scan(&exists)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return false, nil
		default:
			return false, err
		}
	}

	return exists, nil
}

func (r *ProductRepository) List(ctx context.Context, q domain.PaginatedProductsQuery) ([]domain.ProductSummary, domain.Meta, error) {
	var query strings.Builder
	query.WriteString(`
		SELECT 
			p.id, p.name, p.slug, p.price, p.sale_price, p.stock, p.category_id, p.brand_id,
			c.id, c.name, c.slug,
			b.id, b.name, b.slug,
			COUNT(p.id) OVER()
		FROM products p
		LEFT JOIN brands b ON p.brand_id = b.id
		LEFT JOIN categories c ON p.category_id = c.id 
		WHERE p.is_active = true AND (p.name ILIKE '%' || $1 || '%' OR p.description ILIKE '%' || $1 || '%')	
	`)

	params := []any{q.Search}
	paramIndex := len(params) + 1

	if len(q.Categories) > 0 {
		query.WriteString(`
			AND p.category_id in (
				WITH RECURSIVE category_tree AS (
					SELECT id FROM categories
					WHERE slug = ANY($` + strconv.Itoa(paramIndex) + `)
					UNION ALL
					SELECT c.id FROM categories
					JOIN category_tree ct ON c.parent_id = ct.id
				)
				SELECT id FROM category_tree
			)
		`)
		params = append(params, pq.Array(q.Categories))
		paramIndex++
	}

	query.WriteString("GROUP BY p.id, c.id, b.id")

	validSortFields := map[string]string{
		"name":  "p.name",
		"price": "p.price",
		"stock": "p.stock",
	}

	sortField, exists := validSortFields[q.SortField]
	if !exists {
		sortField = "p.name"
	}

	query.WriteString(" ORDER BY ")
	query.WriteString(sortField)
	query.WriteString(" ")
	query.WriteString(q.SortDirection)

	query.WriteString(" LIMIT $")
	query.WriteString(strconv.Itoa(paramIndex))
	params = append(params, q.Limit)
	paramIndex++

	query.WriteString(" OFFSET $")
	query.WriteString(strconv.Itoa(paramIndex))
	params = append(params, q.Offset)

	ctx, cancel := context.WithTimeout(ctx, postgres.QueryTimeoutDuration)
	defer cancel()

	var products []domain.ProductSummary
	var count int
	rows, err := r.db.QueryContext(ctx, query.String(), params...)
	if err != nil {
		return nil, domain.Meta{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var product domain.ProductSummary
		product.Category = &domain.CategorySummary{}
		product.Brand = &domain.BrandSummary{}

		err = rows.Scan(
			&product.Id,
			&product.Name,
			&product.Slug,
			&product.Price,
			&product.SalePrice,
			&product.Stock,
			&product.CategoryId,
			&product.BrandId,
			&product.Category.Id,
			&product.Category.Name,
			&product.Category.Slug,
			&product.Brand.Id,
			&product.Brand.Name,
			&product.Brand.Slug,
			&count,
		)
		if err != nil {
			return nil, domain.Meta{}, err
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, domain.Meta{}, err
	}

	currentPage := (q.Offset / q.Limit) + 1
	totalPages := int(math.Ceil(float64(count) / float64(q.Limit)))
	meta := domain.Meta{
		TotalItems:  count,
		CurrentPage: currentPage,
		PageSize:    q.Limit,
		TotalPages:  totalPages,
	}

	return products, meta, nil
}
