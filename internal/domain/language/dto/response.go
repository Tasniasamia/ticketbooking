package dto;

type Response struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	RTL       bool   `json:"rtl"`
	IsActive  bool   `json:"is_active"`
	IsDefault bool   `json:"is_default"`
	Flag      string `json:"flag,omitempty"`
}