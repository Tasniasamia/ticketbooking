package dto
type CreateTestimonialRequest struct {
	Name string `json:"name" validate:"required,min=2"`
	Role string `json:"role"`
	Text string `json:"text" validate:"required,min=5"`
	ImageURL string `json:"image_url"`
	ImageID uint `json:"image_id"`
	IsActive *bool `json:"is_active"`
	SortOrder *int `json:"sort_order"`
}
type UpdateTestimonialRequest struct {
	Name *string `json:"name" validate:"omitempty,min=2"`
	Role *string `json:"role"`
	Text *string `json:"text" validate:"omitempty,min=5"`
	ImageURL *string `json:"image_url"`
	ImageID *uint `json:"image_id"`
	IsActive *bool `json:"is_active"`
	SortOrder *int `json:"sort_order"`
}
