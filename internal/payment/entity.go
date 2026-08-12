package payment

import (
	"ticketBooking/internal/payment/dto"
	"time"

	"gorm.io/gorm"
)

type PaymentStatus string
type PaymentMethodCode string

const (
	StatusPending   PaymentStatus = "pending"
	StatusSuccess   PaymentStatus = "success"
	StatusFailed    PaymentStatus = "failed"
	StatusCancelled PaymentStatus = "cancelled"

	MethodStripe     PaymentMethodCode = "stripe"
	MethodSSLCommerz PaymentMethodCode = "sslcommerz"
)

// Payment records a single payment attempt for a booking.
type Payment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	BookingID     uint              `gorm:"not null;index" json:"booking_id"`
	UserID        uint              `gorm:"not null;index" json:"user_id"`
	EventID       uint              `gorm:"not null;index" json:"event_id"`
	Amount        float64           `gorm:"not null" json:"amount"`
	Currency      string            `gorm:"size:10;not null" json:"currency"`
	Reason        string            `gorm:"size:255" json:"reason"`
	PaymentMethod PaymentMethodCode `gorm:"size:20;not null" json:"payment_method"`
	Status        PaymentStatus     `gorm:"size:20;not null;default:'pending';index" json:"status"`
	TransactionID string            `gorm:"size:64;uniqueIndex;not null" json:"transaction_id"`

	GatewaySessionID string     `gorm:"size:255;index" json:"gateway_session_id"`
	CheckoutURL      string     `gorm:"type:text" json:"checkout_url"`
	GatewayResponse  string     `gorm:"type:text" json:"-"`
	PaidAt           *time.Time `json:"paid_at"`
}

func (Payment) TableName() string { return "payments" }

func (p *Payment) ToResponse(bookingCode string) *dto.PaymentResponse {
	paidAt := ""
	if p.PaidAt != nil {
		paidAt = p.PaidAt.Format("2006-01-02 15:04:05")
	}
	return &dto.PaymentResponse{
		ID: p.ID, BookingID: p.BookingID, UserID: p.UserID, EventID: p.EventID,
		Amount: p.Amount, Currency: p.Currency, Reason: p.Reason,
		PaymentMethod: string(p.PaymentMethod), Status: string(p.Status),
		TransactionID: p.TransactionID, GatewaySessionID: p.GatewaySessionID,
		CheckoutURL: p.CheckoutURL, BookingCode: bookingCode, PaidAt: paidAt,
		CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// PaymentMethod is a configurable gateway entry (name, logo, enable flag).
// Config credentials still live in settings; only enabled methods are usable at checkout.
type PaymentMethod struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name   string `gorm:"size:100;not null" json:"name"`
	Code   string `gorm:"size:20;uniqueIndex;not null" json:"code"` // stripe | sslcommerz
	LogoURL   string `gorm:"type:text" json:"logo_url"`
	LogoID    uint `gorm:"not null" json:"logo_id"`

	Enable bool   `gorm:"not null;default:true" json:"enable"`
}

func (PaymentMethod) TableName() string { return "payment_methods" }

func (m *PaymentMethod) ToResponse() *dto.PaymentMethodResponse {
	return &dto.PaymentMethodResponse{
		ID:        m.ID,
		Name:      m.Name,
		Code:      m.Code,
		LogoURL:   m.LogoURL,
		LogoID:    m.LogoID,
		Enable:    m.Enable,
		CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: m.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
