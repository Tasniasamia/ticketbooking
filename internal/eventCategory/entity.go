package eventCategory;

import (
	"ticketBooking/internal/eventCategory/dto"
	"ticketBooking/internal/utils/i18n"
	"time"

	"gorm.io/gorm"
)

type EventCategory struct {
	gorm.Model
	Name            i18n.LocalizedString `json:"name" gorm:"type:jsonb;not null"`
	Description      i18n.LocalizedString `json:"description" gorm:"type:jsonb"`
	EventCategoryImageURL    string               `json:"event_category_image_url" gorm:"not null"`
	EventCategoryImageId     int                  `json:"event_category_image_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

}

func (c *EventCategory) ToResponse(lang string) *dto.Response {
	return &dto.Response{
		ID:               c.ID,
		Name:            c.Name.Get(lang),
		Description:      c.Description.Get(lang),
		EventCategoryImageURL:    c.EventCategoryImageURL,
		EventCategoryImageId:     c.EventCategoryImageId,
	    CreatedAt:        c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        c.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (e *EventCategory) ToRawResponse() *dto.RawResponse {
	return &dto.RawResponse{
		ID:               e.ID,
		Name:            e.Name,
		Description:      e.Description,
		EventCategoryImageURL:    e.EventCategoryImageURL,
		EventCategoryImageId:     e.EventCategoryImageId,
		CreatedAt:        e.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        e.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	
}