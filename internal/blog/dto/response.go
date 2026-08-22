package dto

import (
	"ticketBooking/internal/media"
	"ticketBooking/internal/utils/i18n"
	"time"
)

type AuthorInfo struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	ProfileImage   string `json:"profile_image,omitempty"`
	ProfileImageId uint   `json:"profile_image_id,omitempty"`
}

// lang-resolved response (public)
type Response struct {
	ID               uint                 `json:"id"`
	Title            string               `json:"title"`
	ShortDescription string               `json:"short_description"`
	LongDescription  string               `json:"long_description"`
	ThumbnailImage   media.MediaImage     `json:"thumbnail_image"`
	Images           media.MediaImageList `json:"images"`
	AuthorID         uint                 `json:"author_id"`
	Author           *AuthorInfo          `json:"author,omitempty"`
	Status           StatusType           `json:"status"`
	Slug             string               `json:"slug"`
	LikeCount        int64                `json:"like_count"`
	CommentCount     int64                `json:"comment_count"`
	IsLiked          bool                 `json:"is_liked"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
}

// admin / raw — full multi-lang object
type RawResponse struct {
	ID               uint                 `json:"id"`
	Title            i18n.LocalizedString `json:"title"`
	ShortDescription i18n.LocalizedString `json:"short_description"`
	LongDescription  i18n.LocalizedString `json:"long_description"`
	ThumbnailImage   media.MediaImage     `json:"thumbnail_image"`
	Images           media.MediaImageList `json:"images"`
	AuthorID         uint                 `json:"author_id"`
	Author           *AuthorInfo          `json:"author,omitempty"`
	Status           StatusType           `json:"status"`
	Slug             string               `json:"slug"`
	LikeCount        int64                `json:"like_count"`
	CommentCount     int64                `json:"comment_count"`
	IsLiked          bool                 `json:"is_liked"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
}

// Single blog detail — comments tree সহ
type DetailResponse struct {
	*RawResponse
	Comments interface{} `json:"comments"` // nested comment tree
}

type LikeResponse struct {
	Liked     bool  `json:"liked"`
	LikeCount int64 `json:"like_count"`
}
