package domain

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BaseProduct struct {
	Id         int64    `json:"id"`
	Name       string   `json:"name"`
	Slug       string   `json:"slug"`
	Price      float64  `json:"price"`
	SalePrice  *float64 `json:"sale_price"`
	Stock      int64    `json:"stock"`
	CategoryId int64    `json:"category_id"`
	BrandId    int64    `json:"brand_id"`
}

type Product struct {
	BaseProduct
	Description *string   `json:"description"`
	Category    *Category `json:"category"`
	Brand       *Brand    `json:"brand"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductSummary struct {
	BaseProduct
	Category *CategorySummary `json:"category"`
	Brand    *BrandSummary    `json:"brand"`
}

type PaginatedProductsQuery struct {
	Offset        int      `json:"offset"`
	Limit         int      `json:"limit"`
	Search        string   `json:"search"`
	SortDirection string   `json:"sort_direction" validate:"oneof=asc desc"`
	SortField     string   `json:"sort_field" validate:"oneof=price name stock"`
	Categories    []string `json:"categories"`
}

func (q PaginatedProductsQuery) Parse(r *http.Request) (PaginatedProductsQuery, error) {
	qs := r.URL.Query()

	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return q, err
		}
		q.Offset = o
	}

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return q, err
		}
		q.Limit = l
	}

	search := qs.Get("search")
	if search != "" {
		q.Search = search
	}

	sort := qs.Get("sort_direction")
	if sort != "" {
		q.SortDirection = sort
	}

	sortField := qs.Get("sort_field")
	if sortField != "" {
		q.SortField = sortField
	}

	categories := qs.Get("categories")
	if categories != "" {
		q.Categories = strings.Split(categories, ",")
	} else {
		q.Categories = []string{}
	}

	return q, nil
}
