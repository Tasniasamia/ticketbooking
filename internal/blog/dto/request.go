package dto

import "ticketBooking/internal/media"

type StatusType string

const (
	Pending  StatusType = "pending"
	Approved StatusType = "approved"
	Canceled StatusType = "canceled"
)

type CreateRequest struct {
	Title            map[string]string    `json:"title" validate:"required"`
	ShortDescription map[string]string    `json:"short_description"`
	LongDescription  map[string]string    `json:"long_description"`
	ThumbnailImage   media.MediaImage     `json:"thumbnail_image" validate:"required"`
	Images           media.MediaImageList `json:"images"`
	AuthorID         uint                 `json:"author_id"`
	Status           StatusType           `json:"status"`
	Slug             string               `json:"slug"`
}

type UpdateRequest struct {
	Title            map[string]string    `json:"title"`
	ShortDescription map[string]string    `json:"short_description"`
	LongDescription  map[string]string    `json:"long_description"`
	ThumbnailImage   media.MediaImage     `json:"thumbnail_image"`
	Images           media.MediaImageList `json:"images"`
	AuthorID         uint                 `json:"author_id"`
	Status           StatusType           `json:"status"`
	Slug             string               `json:"slug"`
}

type UpdateBlogStatusRequest struct {
	BlogID uint       `json:"blogId" validate:"required"`
	Status StatusType `json:"status" validate:"required"`
}
