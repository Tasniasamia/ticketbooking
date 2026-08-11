package dto;

type CreateCurrencyRequest struct {
	Code      string  `json:"code" validate:"required,min=3,max=10"`
	Name      string  `json:"name" validate:"required"`
	Symbol    string  `json:"symbol"`
	RateToBDT float64 `json:"rate_to_bdt" validate:"required,gt=0"`
	Status    string  `json:"status"` // enable | disable
	IsDefault bool    `json:"is_default"`
}

type UpdateCurrencyRequest struct {
	Name      *string  `json:"name"`
	Symbol    *string  `json:"symbol"`
	RateToBDT *float64 `json:"rate_to_bdt"`
	Status    *string  `json:"status"`
}

type ConvertRequest struct {
	Amount   float64 `json:"amount" validate:"required,gt=0"`
	FromCode string  `json:"from_code" validate:"required"`
	ToCode   string  `json:"to_code" validate:"required"`
}



type SetDefaultRequest struct {
	Code string `json:"code" validate:"required"`
}
