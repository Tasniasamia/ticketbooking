package dto

import "time"

type CreateRequest struct {
	Title        map[string]string `json:"title" validate:"required"`
	Description  map[string]string `json:"description"`
	Location     map[string]string `json:"location" validate:"required"`
	StartsAt     time.Time         `json:"starts_at" validate:"required"`
	TotalTickets int               `json:"total_tickets" validate:"required,gt=0"`
	Price        int               `json:"price" validate:"gte=0"`
	EventImageURL    string               `json:"event_url" validate:"required"`
	EventImageId     int                  `json:"event_image_id" validate:"required"`
}

type UpdateRequest struct {
	Title        map[string]string `json:"title"`
	Description  map[string]string `json:"description"`
	Location     map[string]string `json:"location"`
	StartsAt     *time.Time        `json:"starts_at"`
	Price        *int              `json:"price" validate:"omitempty,gte=0"`
	EventImageURL    *string               `json:"event_url"`
	EventImageId     *int                  `json:"event_image_id"`
}