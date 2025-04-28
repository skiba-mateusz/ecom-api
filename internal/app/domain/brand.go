package domain

type Brand struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	LogoUrl     *string `json:"logo_url"`
}
