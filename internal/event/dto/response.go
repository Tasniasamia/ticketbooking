package dto

import (
	"ticketBooking/internal/utils/i18n"
	"time"
)

// // lang resolve করে single string
type Response struct {
	ID               uint      `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Location         string    `json:"location"`
	StartsAt         time.Time `json:"starts_at"`
	TotalTickets     int       `json:"total_tickets"`
	AvailableTickets int       `json:"available_tickets"`
	Price            int       `json:"price"`
	CreatedAt        string    `json:"created_at"`
	EventImageURL    string    `json:"event_url"`
	EventImageId     int       `json:"event_image_id"`
}

// admin / raw — পুরো multi-lang object
type RawResponse struct {
	ID               uint                 `json:"id"`
	Title            i18n.LocalizedString `json:"title"`
	Description      i18n.LocalizedString `json:"description"`
	Location         i18n.LocalizedString `json:"location"`
	StartsAt         time.Time            `json:"starts_at"`
	TotalTickets     int                  `json:"total_tickets"`
	AvailableTickets int                  `json:"available_tickets"`
	Price            int                  `json:"price"`
	CreatedAt        string               `json:"created_at"`
	EventImageURL    string               `json:"event_url"`
	EventImageId     int                  `json:"event_image_id"`
}