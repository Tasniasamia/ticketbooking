package dto

import (
	"ticketBooking/internal/media"
	"ticketBooking/internal/user/dto"

	"ticketBooking/internal/utils/i18n"
	"time"
)

type ManagerInfo struct {
	ID             uint         `json:"id"`
	Name           string       `json:"name"`
	Email          string       `json:"email"`
	Role           dto.RoleType `json:"role"`
	Designation    string       `json:"designation,omitempty"`
	ProfileImage   string       `json:"profile_image,omitempty"`
	ProfileImageId uint         `json:"profile_image_id,omitempty"`
	PhoneNumber    string       `json:"phone_number,omitempty"`
	Country        string       `json:"country,omitempty"`
	Status         string       `json:"status,omitempty"`
}

type CategoryInfo struct {
	ID                    uint                 `json:"id"`
	Name                  i18n.LocalizedString `json:"name" gorm:"type:jsonb;not null"`
	Description           i18n.LocalizedString `json:"description" gorm:"type:jsonb"`
	EventCategoryImageURL string               `json:"event_category_image_url" gorm:"not null"`
	EventCategoryImageId  int                  `json:"event_category_image_id" gorm:"not null"`
}

// // lang resolve করে single string
type Response struct {
	ID               uint                 `json:"id"`
	Title            string               `json:"title"`
	Description      string               `json:"description"`
	Location         string               `json:"location"`
	StartsAt         time.Time            `json:"starts_at"`
	TotalTickets     int                  `json:"total_tickets"`
	AvailableTickets int                  `json:"available_tickets"`
	Price            int                  `json:"price"`
	CreatedAt        time.Time            `json:"created_at"`
	ThumbnailImage   media.MediaImage     `json:"thumbnail_image" validate:"required"`
	Images           media.MediaImageList `json:"images"`
	ManagerID        uint                 `json:"manager_id"`
	Manager          *ManagerInfo         `json:"manager,omitempty"` // ← এখানে

	CategoryID uint          `json:"category_id"`
	Category   *CategoryInfo `json:"category,omitempty"` // ← এখানেeventCategory.EventCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

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
	CreatedAt        time.Time            `json:"created_at"`
	ThumbnailImage   media.MediaImage     `json:"thumbnail_image" validate:"required"`
	Images           media.MediaImageList `json:"images"`
	ManagerID        uint                 `json:"manager_id"`
	Manager          *ManagerInfo         `json:"manager,omitempty"` // ← এখানে

	CategoryID uint          `json:"category_id"`
	Category   *CategoryInfo `json:"category,omitempty"` // ← এখানে
}
