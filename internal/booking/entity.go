package booking

import (
	"ticketBooking/internal/booking/dto"
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	ID          uint              `gorm:"primaryKey" json:"id"`
	UserID      uint              `gorm:"not null;index" json:"user_id"`
	EventID     uint              `gorm:"not null;index" json:"event_id"`
	Quantity    int               `gorm:"not null" json:"quantity"`
	TotalPrice  float64           `gorm:"not null" json:"total_price"`
	Status      dto.BookingStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	BookingCode string            `gorm:"uniqueIndex;size:50" json:"booking_code"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"-"`
    

}

func (Booking) TableName() string { return "bookings" }

func (b *Booking) ToResponse() *dto.Response {
	return &dto.Response{
		ID:          b.ID,
		UserID:      b.UserID,
		EventID:     b.EventID,
		Quantity:    b.Quantity,
		TotalPrice:  b.TotalPrice,
		Status:      b.Status,
		BookingCode: b.BookingCode,
		CreatedAt:   b.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
