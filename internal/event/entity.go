package event

import (
	"ticketBooking/internal/event/dto"
	"ticketBooking/internal/utils/i18n"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title            i18n.LocalizedString `json:"title" gorm:"type:jsonb;not null"`
	Description      i18n.LocalizedString `json:"description" gorm:"type:jsonb"`
	Location         i18n.LocalizedString `json:"location" gorm:"type:jsonb;not null"`
	StartsAt         time.Time            `json:"starts_at" gorm:"not null"`
	TotalTickets     int                  `json:"total_tickets" gorm:"not null"`
	AvailableTickets int                  `json:"available_tickets" gorm:"not null"`
	Price            int                  `json:"price" gorm:"not null"`
}

func (e *Event) ToResponse(lang string) *dto.Response {
	return &dto.Response{
		ID:               e.ID,
		Title:            e.Title.Get(lang),
		Description:      e.Description.Get(lang),
		Location:         e.Location.Get(lang),
		StartsAt:         e.StartsAt,
		TotalTickets:     e.TotalTickets,
		AvailableTickets: e.AvailableTickets,
		Price:            e.Price,
		CreatedAt:        e.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (e *Event) ToRawResponse() *dto.RawResponse {
	return &dto.RawResponse{
		ID:               e.ID,
		Title:            e.Title,
		Description:      e.Description,
		Location:         e.Location,
		StartsAt:         e.StartsAt,
		TotalTickets:     e.TotalTickets,
		AvailableTickets: e.AvailableTickets,
		Price:            e.Price,
		CreatedAt:        e.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}