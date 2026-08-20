package dto

import (
	bookingDto "ticketBooking/internal/booking/dto"
	eventDto "ticketBooking/internal/event/dto"
	"ticketBooking/internal/media"
	userDto "ticketBooking/internal/user/dto"
	"ticketBooking/internal/utils/i18n"

	"time"
)

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
type UserInfo struct {
	ID             uint         `json:"id"`
	Name           string       `json:"name"`
	Email          string       `json:"email"`
	Role           userDto.RoleType `json:"role"`
	ProfileImage   string       `json:"profile_image,omitempty"`
	ProfileImageId uint         `json:"profile_image_id,omitempty"`
	PhoneNumber    string       `json:"phone_number,omitempty"`
	Country        string       `json:"country,omitempty"`
	Status         string       `json:"status,omitempty"`
}

type BookingInfo struct{
	ID          uint              `gorm:"primaryKey" json:"id"`
	UserID      uint              `gorm:"not null;index" json:"user_id"`
	EventID     uint              `gorm:"not null;index" json:"event_id"`
	Quantity    int               `gorm:"not null" json:"quantity"`
	TotalPrice  float64           `gorm:"not null" json:"total_price"`
	Status      bookingDto.BookingStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	BookingCode string            `gorm:"uniqueIndex;size:50" json:"booking_code"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type EventInfo struct {
	ID               uint                 `json:"id"`
	Title            i18n.LocalizedString               `json:"title"`
	Description      i18n.LocalizedString               `json:"description"`
	Location         i18n.LocalizedString               `json:"location"`
	StartsAt         time.Time            `json:"starts_at"`
	TotalTickets     int                  `json:"total_tickets"`
	AvailableTickets int                  `json:"available_tickets"`
	Price            int                  `json:"price"`
	Currency         string               `json:"currency"`
	CreatedAt        time.Time            `json:"created_at"`
	ThumbnailImage   media.MediaImage     `json:"thumbnail_image" validate:"required"`
	Images           media.MediaImageList `json:"images"`
	ManagerID        uint                 `json:"manager_id"`
	Manager          *UserInfo              `json:"manager,omitempty"` // ← এখানে
	Status         eventDto.StatusType       `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
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
	Bookings         BookingInfo    `json:"bookings"`

	UserInfo       UserInfo       `json:"user_info"`
    
	EventInfo     EventInfo     `json:"event_info"`
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

// PaymentSummary — amounts converted to site default currency.
// pending = status pending sum, paid = status success sum, total = pending + paid.
type PaymentSummary struct {
	Pending  float64 `json:"pending"`
	Paid     float64 `json:"paid"`
	Total    float64 `json:"total"`
	Currency string  `json:"currency"` // default currency code
}
