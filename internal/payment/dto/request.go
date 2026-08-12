package dto

type CreateCheckoutRequest struct {
	EventID       uint   `json:"event_id" validate:"required"`
	Quantity      int    `json:"quantity" validate:"required,min=1"`
	PaymentMethod string `json:"payment_method" validate:"required,oneof=stripe sslcommerz"`   //PaymentMethodCode

	// Customer billing details (used mainly by SSLCommerz; optional for Stripe)
	CustomerName        string `json:"customer_name" validate:"required"`
	CustomerPhoneNumber string `json:"customer_phone_number" validate:"required"`
	CustomerEmail       string `json:"customer_email" validate:"required"`
	CustomerAddress     string `json:"customer_address" validate:"required"`
	Country             string `json:"country" validate:"required"`
	Postcode            string `json:"postcode" validate:"required"`
}

type SSLCommerzIPN struct {
	TranID     string `form:"tran_id" json:"tran_id"`
	ValID      string `form:"val_id" json:"val_id"`
	Amount     string `form:"amount" json:"amount"`
	Status     string `form:"status" json:"status"`
	BankTranID string `form:"bank_tran_id" json:"bank_tran_id"`
	Currency   string `form:"currency" json:"currency"`
	CardType   string `form:"card_type" json:"card_type"`
	ValueA     string `form:"value_a" json:"value_a"`
	ValueB     string `form:"value_b" json:"value_b"`
	ValueC     string `form:"value_c" json:"value_c"`
	ValueD     string `form:"value_d" json:"value_d"`
}

// ---- Payment Method CRUD ----

type CreatePaymentMethodRequest struct {
	Name   string `json:"name" validate:"required,min=1,max=100"`
	Code   string `json:"code" validate:"required,oneof=stripe sslcommerz"`
	LogoURL   string `json:"logo_url" validate:"required"`
	LogoID    string `json:"logo_id" validate:"required"`
	Enable *bool  `json:"enable"`
}

type UpdatePaymentMethodRequest struct {
	Name   *string `json:"name" validate:"omitempty,min=1,max=100"`
	Code   *string `json:"code" validate:"omitempty,oneof=stripe sslcommerz"`
	LogoUrl   *string `json:"logo_url" validate:"omitempty"`
	LogoID    *string `json:"logo_id" validate:"omitempty"`
	Enable *bool   `json:"enable"`
}
