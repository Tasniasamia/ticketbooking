package payment

import (
	"ticketBooking/internal/booking"
	"ticketBooking/internal/event"
	"ticketBooking/internal/payment/dto"
	"ticketBooking/internal/user"
	"time"

	// "gorm.io/datatypes"

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

	BookingID uint            `gorm:"not null;index" json:"booking_id"`
	Bookings  booking.Booking `json:"bookings,omitempty" gorm:"foreignKey:BookingID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	UserID   uint      `gorm:"not null;index" json:"user_id"`
	UserInfo user.User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	EventID       uint              `gorm:"not null;index" json:"event_id"`
	EventInfo     event.Event       `json:"event,omitempty" gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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
		Bookings: dto.BookingInfo{
			ID:          p.Bookings.ID,
			BookingCode: p.Bookings.BookingCode,
			UserID:      p.UserID,
			EventID:     p.EventID,
			Quantity:    p.Bookings.Quantity,
			TotalPrice:  p.Bookings.TotalPrice,
			Status:      p.Bookings.Status,
		},
		UserInfo: dto.UserInfo{
			ID:             p.UserInfo.ID,
			Name:           p.UserInfo.Name,
			Email:          p.UserInfo.Email,
			Role:           p.UserInfo.Role,
			ProfileImage:   p.UserInfo.ProfileImage,
			ProfileImageId: p.UserInfo.ProfileImageId,
			PhoneNumber:    p.UserInfo.PhoneNumber,
			Country:        p.UserInfo.Country,
			Status:         string(p.UserInfo.Status),
		},
		EventInfo: dto.EventInfo{
			ID:               p.EventInfo.ID,
			Title:            p.EventInfo.Title,
			Description:      p.EventInfo.Description,
			Location:         p.EventInfo.Location,
			StartsAt:         p.EventInfo.StartsAt,
			TotalTickets:     p.EventInfo.TotalTickets,
			AvailableTickets: p.EventInfo.AvailableTickets,
			Price:            p.EventInfo.Price,
			Currency:         p.EventInfo.Currency,
			CreatedAt:        p.EventInfo.CreatedAt,
			ThumbnailImage:   p.EventInfo.ThumbnailImage,
			Images:           p.EventInfo.Images,
			ManagerID:        p.EventInfo.ManagerID,
			Manager: &dto.UserInfo{
				ID:             p.UserInfo.ID,
				Name:           p.UserInfo.Name,
				Email:          p.UserInfo.Email,
				Role:           p.UserInfo.Role,
				ProfileImage:   p.UserInfo.ProfileImage,
				ProfileImageId: p.UserInfo.ProfileImageId,
				PhoneNumber:    p.UserInfo.PhoneNumber,
				Country:        p.UserInfo.Country,
				Status:         string(p.UserInfo.Status),
			},

			Status: p.EventInfo.Status,
		},
	}
}

func (p *Payment) ToRawResponse() *dto.PaymentResponse {
	paidAt := ""
	if p.PaidAt != nil {
		paidAt = p.PaidAt.Format("2006-01-02 15:04:05")
	}
	return &dto.PaymentResponse{
		ID: p.ID, BookingID: p.BookingID, UserID: p.UserID, EventID: p.EventID,
		Amount: p.Amount, Currency: p.Currency, Reason: p.Reason,
		PaymentMethod: string(p.PaymentMethod), Status: string(p.Status),
		TransactionID: p.TransactionID, GatewaySessionID: p.GatewaySessionID,
		CheckoutURL: p.CheckoutURL,
		BookingCode: p.Bookings.BookingCode,
		PaidAt:      paidAt,
		CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02 15:04:05"),
		Bookings: dto.BookingInfo{
			ID:          p.Bookings.ID,
			BookingCode: p.Bookings.BookingCode,
			UserID:      p.UserID,
			EventID:     p.EventID,
			Quantity:    p.Bookings.Quantity,
			TotalPrice:  p.Bookings.TotalPrice,
			Status:      p.Bookings.Status,
		},
		UserInfo: dto.UserInfo{
			ID:             p.UserInfo.ID,
			Name:           p.UserInfo.Name,
			Email:          p.UserInfo.Email,
			Role:           p.UserInfo.Role,
			ProfileImage:   p.UserInfo.ProfileImage,
			ProfileImageId: p.UserInfo.ProfileImageId,
			PhoneNumber:    p.UserInfo.PhoneNumber,
			Country:        p.UserInfo.Country,
			Status:         string(p.UserInfo.Status),
		},
		EventInfo: dto.EventInfo{
			ID:               p.EventInfo.ID,
			Title:            p.EventInfo.Title,
			Description:      p.EventInfo.Description,
			Location:         p.EventInfo.Location,
			StartsAt:         p.EventInfo.StartsAt,
			TotalTickets:     p.EventInfo.TotalTickets,
			AvailableTickets: p.EventInfo.AvailableTickets,
			Price:            p.EventInfo.Price,
			Currency:         p.EventInfo.Currency,
			CreatedAt:        p.EventInfo.CreatedAt,
			ThumbnailImage:   p.EventInfo.ThumbnailImage,
			Images:           p.EventInfo.Images,
			ManagerID:        p.EventInfo.ManagerID,
			Manager: &dto.UserInfo{
				ID:             p.UserInfo.ID,
				Name:           p.UserInfo.Name,
				Email:          p.UserInfo.Email,
				Role:           p.UserInfo.Role,
				ProfileImage:   p.UserInfo.ProfileImage,
				ProfileImageId: p.UserInfo.ProfileImageId,
				PhoneNumber:    p.UserInfo.PhoneNumber,
				Country:        p.UserInfo.Country,
				Status:         string(p.UserInfo.Status),
			},

			Status: p.EventInfo.Status,
		},
	}
}

// PaymentMethod is a configurable gateway entry (name, logo, enable flag).
// Config credentials still live in settings; only enabled methods are usable at checkout.

type PaymentMethod struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name        string `gorm:"size:100;not null"`
	Code        string `gorm:"size:20;uniqueIndex;not null"`
	LogoURL     string `gorm:"type:text"`
	LogoID      uint   `gorm:"not null"`
	Enable      bool   `gorm:"not null;default:true"`
	Credentials string `gorm:"type:text"`
}

func (PaymentMethod) TableName() string { return "payment_methods" }

func (m *PaymentMethod) ToResponse() *dto.PaymentMethodResponse {

	resp := &dto.PaymentMethodResponse{
		ID:        m.ID,
		Name:      m.Name,
		Code:      m.Code,
		LogoURL:   m.LogoURL,
		LogoID:    m.LogoID,
		Enable:    m.Enable,
		CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: m.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if m.Credentials != "" {
		if creds, err := DecryptCredentials(m.Credentials); err == nil {
			resp.Credentials = creds
		}
	}
	return resp
}
