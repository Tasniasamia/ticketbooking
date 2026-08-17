package dto

import (
	"ticketBooking/internal/media"
	"time"
)



type CreateRequest struct {
	Title          map[string]string `json:"title" validate:"required"`
	Description    map[string]string `json:"description"`
	Location       map[string]string `json:"location" validate:"required"`
	StartsAt       time.Time         `json:"starts_at" validate:"required"`
	TotalTickets   int               `json:"total_tickets" validate:"required,gt=0"`
	Price          int               `json:"price" validate:"gte=0"`
	ThumbnailImage media.MediaImage     `json:"thumbnail_image" validate:"required"`
	Images         media.MediaImageList `json:"images"`
	ManagerID      uint              `json:"manager_id" gorm:"not null;index"`
	CategoryID     uint              `json:"category_id" gorm:"not null;index"`
}

type UpdateRequest struct {
	Title          map[string]string `json:"title"`
	Description    map[string]string `json:"description"`
	Location       map[string]string `json:"location"`
	StartsAt       time.Time         `json:"starts_at"`
	Price          int               `json:"price" validate:"omitempty,gte=0"`
	ThumbnailImage media.MediaImage     `json:"thumbnail_image" validate:"required"`
	Images         media.MediaImageList `json:"images"`
	ManagerID      uint              `json:"manager_id" gorm:"not null;index"`
	CategoryID     uint              `json:"category_id" gorm:"not null;index"`
}

