package dto;

type Response struct {
	ID               uint      `json:"id"`
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Symbol    string  `json:"symbol"`
	RateToBDT float64 `json:"rate_to_bdt"`
	Status    string  `json:"status"`
	IsDefault bool    `json:"is_default"`
}

type ConvertResponse struct {
	Amount   float64 `json:"amount"`
	FromCode string  `json:"from_code"`
	ToCode   string  `json:"to_code"`
	Result   float64 `json:"result"`
}