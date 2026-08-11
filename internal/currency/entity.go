package currency

import (
	"time"
	"ticketBooking/internal/currency/dto"

	"gorm.io/gorm"
)

// Currency holds supported currencies. Base currency is always BDT.
type Currency struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Code      string         `gorm:"size:10;uniqueIndex;not null" json:"code"` // BDT, USD, INR, EUR...
	Name      string         `gorm:"size:100;not null" json:"name"`
	Symbol    string         `gorm:"size:10" json:"symbol"`                   // ৳, $, ₹
	RateToBDT float64        `gorm:"not null;default:1" json:"rate_to_bdt"`   // 1 unit of this currency = ? BDT
	IsDefault bool           `gorm:"default:false" json:"is_default"`
	Status    string         `gorm:"size:20;default:'enable'" json:"status"` // enable / disable
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}


func (c *Currency) ToResponse() *dto.Response {
	return &dto.Response{
		ID:               c.ID,
		Code:             c.Code,
		Name:             c.Name,
		Symbol:           c.Symbol,
		RateToBDT:        c.RateToBDT,
		Status:           c.Status,
		IsDefault:        c.IsDefault,
	}

}