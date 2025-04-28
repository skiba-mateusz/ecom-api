package domain

import "time"

type Product struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description"`
	Price       float64   `json:"price"`
	SalePrice   *float64  `json:"sale_price"`
	Stock       int64     `json:"stock"`
	CategoryId  int64     `json:"category_id"`
	Category    *Category `json:"category"`
	BrandId     int64     `json:"brand_id"`
	Brand       *Brand    `json:"brand"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
