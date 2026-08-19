package dto

type CheckoutResponse struct {
	PaymentID     uint    `json:"payment_id"`
	BookingID     uint    `json:"booking_id"`
	BookingCode   string  `json:"booking_code"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Reason        string  `json:"reason"`
	PaymentMethod string  `json:"payment_method"`
	TransactionID string  `json:"transaction_id"`
	CheckoutURL   string  `json:"checkout_url"`
	SessionID     string  `json:"session_id,omitempty"`
	Status        string  `json:"status"`
	CreatedAt     string  `json:"created_at"`
}

type PaymentResponse struct {
	ID               uint    `json:"id"`
	BookingID        uint    `json:"booking_id"`
	UserID           uint    `json:"user_id"`
	EventID          uint    `json:"event_id"`
	Amount           float64 `json:"amount"`
	Currency         string  `json:"currency"`
	Reason           string  `json:"reason"`
	PaymentMethod    string  `json:"payment_method"`
	Status           string  `json:"status"`
	TransactionID    string  `json:"transaction_id"`
	GatewaySessionID string  `json:"gateway_session_id,omitempty"`
	CheckoutURL      string  `json:"checkout_url,omitempty"`
	BookingCode      string  `json:"booking_code,omitempty"`
	PaidAt           string  `json:"paid_at,omitempty"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

type PaymentMethodResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	LogoURL   string `json:"logo_url"`
	LogoID    uint   `json:"logo_id"`
	Enable    bool   `json:"enable"`
   Credentials map[string]string `json:"credentials" `
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
