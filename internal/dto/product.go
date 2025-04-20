package dto

type Product struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryId  uint    `json:"category_id"`
	ImageUrl    string  `json:"image_url"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type UpdateStock struct {
	Stock int `json:"stock"`
}

type UpdateProduct struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	CategoryId  *uint    `json:"category_id"`
	ImageUrl    *string  `json:"image_url"`
	Price       *float64 `json:"price"`
	Stock       *int     `json:"stock"`
}
