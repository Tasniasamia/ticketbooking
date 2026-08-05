package dto

type CreateRequest struct {
	Name      string `json:"name" validate:"required,min=2,max=50"`
	Code      string `json:"code" validate:"required,min=2,max=5"`
	RTL       bool   `json:"rtl"`
	IsActive  *bool  `json:"is_active"`
	IsDefault bool   `json:"is_default"`
	Flag      string `json:"flag"`
}

type UpdateRequest struct {
	Name      string `json:"name" validate:"omitempty,min=2,max=50"`
	Code      string `json:"code" validate:"omitempty,min=2,max=5"`
	RTL       *bool  `json:"rtl"`
	IsActive  *bool  `json:"is_active"`
	IsDefault *bool  `json:"is_default"`
	Flag      string `json:"flag"`
}