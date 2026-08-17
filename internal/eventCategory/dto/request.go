package dto

import (
	"time"
)

type CreateEventCategoryRequest struct {
	Name             map[string]string`json:"name" gorm:"type:jsonb;not null"`
	Description       map[string]string`json:"description" gorm:"type:jsonb"`
	EventCategoryImageURL    string               `json:"event_category_image_url" gorm:"not null"`
	EventCategoryImageId     int                  `json:"event_category_image_id" gorm:"not null"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type UpdateEventCategoryRequest struct {
	Name            map[string]string`json:"name" gorm:"type:jsonb;not null"`
	Description       map[string]string`json:"description" gorm:"type:jsonb"`
	EventCategoryImageURL    string               `json:"event_category_image_url" gorm:"not null"`
	EventCategoryImageId     int                  `json:"event_category_image_id" gorm:"not null"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}