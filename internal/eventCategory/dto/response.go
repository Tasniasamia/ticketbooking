package dto

import (
	"ticketBooking/internal/utils/i18n"
	
)

type Response struct {
	ID                    uint                 `json:"id"`
	Name                  string `json:"name" gorm:"type:jsonb;not null"`
	Description           string `json:"description" gorm:"type:jsonb"`
	EventCategoryImageURL string               `json:"event_category_image_url" gorm:"not null"`
	EventCategoryImageId  int                  `json:"event_category_image_id" gorm:"not null"`
	CreatedAt             string          `json:"created_at"`
	UpdatedAt             string           `json:"updated_at"`
}

// admin / raw — পুরো multi-lang object
type RawResponse struct {
	ID                    uint                 `json:"id"`
	Name                  i18n.LocalizedString `json:"name" gorm:"type:jsonb;not null"`
	Description           i18n.LocalizedString `json:"description" gorm:"type:jsonb"`
	EventCategoryImageURL string               `json:"event_category_image_url" gorm:"not null"`
	EventCategoryImageId int                  `json:"event_category_image_id" gorm:"not null"`
	CreatedAt             string          `json:"created_at"`
	UpdatedAt             string           `json:"updated_at"`
}
