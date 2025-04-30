package domain

type Meta struct {
	TotalItems  int `json:"total_items"`
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalPages  int `json:"total_pages"`
}
