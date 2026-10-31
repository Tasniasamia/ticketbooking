package dto

type CreateFAQRequest struct {
	Question  string `json:"question" validate:"required,min=3"`
	Answer    string `json:"answer" validate:"required,min=3"`
	IsActive  *bool  `json:"is_active"`
	SortOrder *int   `json:"sort_order"`
}
type UpdateFAQRequest struct {
	Question  *string `json:"question" validate:"omitempty,min=3"`
	Answer    *string `json:"answer" validate:"omitempty,min=3"`
	IsActive  *bool   `json:"is_active"`
	SortOrder *int    `json:"sort_order"`
}
