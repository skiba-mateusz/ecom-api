package domain

type Category struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description"`
	ParentId    *int64    `json:"parent_id"`
	Parent      *Category `json:"parent"`
	ImageUrl    *string   `json:"image_url"`
}

type CategorySummary struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
