package dto

type Category struct {
	Name         string `json:"name"`
	ParentId     uint   `json:"parent_id"`
	ImageUrl     string `json:"image_url"`
	DisplayOrder int    `json:"display_order"`
}

type UpdateCategory struct {
	Name         *string `json:"name"`
	ParentId     *uint   `json:"parent_id"`
	ImageUrl     *string `json:"image_url"`
	DisplayOrder *int    `json:"display_order"`
}
