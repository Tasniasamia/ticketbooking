package event

import (
	"ticketBooking/internal/event/dto"
	"ticketBooking/internal/eventCategory"
	"ticketBooking/internal/media"
	"ticketBooking/internal/user"
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

	// --- Manager relation ---
	ManagerID uint      `json:"manager_id" gorm:"not null;index"`
	Manager   user.User `json:"manager,omitempty" gorm:"foreignKey:ManagerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// --- Images (new structure) ---
	ThumbnailImage media.MediaImage     `json:"thumbnail_image" validate:"required"`
	Images         media.MediaImageList `json:"images"`   // array of objects
	CategoryID uint                       `json:"category_id" gorm:"not null;index"`
	Category   eventCategory.EventCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Status         dto.StatusType       `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`

}


func (e *Event) ToResponse() *dto.RawResponse {
	res := &dto.RawResponse{
		ID:               e.ID,
		Title:            e.Title,
		Description:      e.Description,
		Location:         e.Location,
		StartsAt:         e.StartsAt,
		TotalTickets:     e.TotalTickets,
		AvailableTickets: e.AvailableTickets,
		Price:            e.Price,
		CreatedAt:        e.CreatedAt,
		ManagerID:        e.ManagerID,
		CategoryID:       e.CategoryID,
		ThumbnailImage:   e.ThumbnailImage,
		Images:           e.Images,
		Status:           e.Status,
	}

	// শুধু যে field চাও, সেগুলোই দাও
	if e.Manager.ID != 0 {
		res.Manager = &dto.ManagerInfo{
			ID:             e.Manager.ID,
			Name:           e.Manager.Name,
			Email:          e.Manager.Email,
			Role:           e.Manager.Role,
			Designation:    e.Manager.Designation,
			ProfileImage:   e.Manager.ProfileImage,
			ProfileImageId: e.Manager.ProfileImageId,
			PhoneNumber:    e.Manager.PhoneNumber,
			Country:        e.Manager.Country,
			Status:         string(e.Manager.Status),
		}
	}

	// Category-ও same ভাবে
	if e.Category.ID != 0 {
		res.Category = &dto.CategoryInfo{
			ID:          e.Category.ID,
			Name:        e.Category.Name,
			Description: e.Category.Description,
			EventCategoryImageURL:    e.Category.EventCategoryImageURL,
			EventCategoryImageId:     e.Category.EventCategoryImageId,
		}
	}

	return res
}
func (e *Event) ToRawResponse() *dto.RawResponse {
	res := &dto.RawResponse{
		ID:               e.ID,
		Title:            e.Title,
		Description:      e.Description,
		Location:         e.Location,
		StartsAt:         e.StartsAt,
		TotalTickets:     e.TotalTickets,
		AvailableTickets: e.AvailableTickets,
		Price:            e.Price,
		CreatedAt:        e.CreatedAt,
		ManagerID:        e.ManagerID,
		CategoryID:       e.CategoryID,
		ThumbnailImage:   e.ThumbnailImage,
		Images:           e.Images,
		Status:           e.Status,

	}

	// শুধু যে field চাও, সেগুলোই দাও
	if e.Manager.ID != 0 {
		res.Manager = &dto.ManagerInfo{
			ID:             e.Manager.ID,
			Name:           e.Manager.Name,
			Email:          e.Manager.Email,
			Role:           e.Manager.Role,
			Designation:    e.Manager.Designation,
			ProfileImage:   e.Manager.ProfileImage,
			ProfileImageId: e.Manager.ProfileImageId,
			PhoneNumber:    e.Manager.PhoneNumber,
			Country:        e.Manager.Country,
			Status:         string(e.Manager.Status),
		}
	}

	// Category-ও same ভাবে
	if e.Category.ID != 0 {
		res.Category = &dto.CategoryInfo{
			ID:          e.Category.ID,
			Name:        e.Category.Name,
			Description: e.Category.Description,
			EventCategoryImageURL:    e.Category.EventCategoryImageURL,
			EventCategoryImageId:     e.Category.EventCategoryImageId,
		}
	}

	return res
}