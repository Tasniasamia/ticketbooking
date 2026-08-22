package blog

import (
	"ticketBooking/internal/blog/dto"
	"ticketBooking/internal/media"
	"ticketBooking/internal/user"
	"ticketBooking/internal/utils/i18n"

	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model

	Title            i18n.LocalizedString `json:"title" gorm:"type:jsonb;not null"`
	ShortDescription i18n.LocalizedString `json:"short_description" gorm:"type:jsonb"`
	LongDescription  i18n.LocalizedString `json:"long_description" gorm:"type:jsonb"`
	ThumbnailImage   media.MediaImage     `json:"thumbnail_image"`
	Images           media.MediaImageList `json:"images"`

	// Author
	AuthorID uint      `json:"author_id" gorm:"not null;index"`
	Author   user.User `json:"author,omitempty" gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Status dto.StatusType `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	Slug   string         `json:"slug" gorm:"type:varchar(255);index"`

	// Denormalized counters — atomic SQL update দিয়ে concurrent-safe
	LikeCount    int64 `json:"like_count" gorm:"not null;default:0"`
	CommentCount int64 `json:"comment_count" gorm:"not null;default:0"`
}

// BlogLike — এক user এক blog-এ একবারই like করতে পারবে (unique index)
type BlogLike struct {
	ID     uint `gorm:"primaryKey"`
	BlogID uint `json:"blog_id" gorm:"not null;uniqueIndex:idx_blog_user_like"`
	UserID uint `json:"user_id" gorm:"not null;uniqueIndex:idx_blog_user_like"`
}

func (BlogLike) TableName() string {
	return "blog_likes"
}

func (b *Blog) ToRawResponse() *dto.RawResponse {
	res := &dto.RawResponse{
		ID:               b.ID,
		Title:            b.Title,
		ShortDescription: b.ShortDescription,
		LongDescription:  b.LongDescription,
		ThumbnailImage:   b.ThumbnailImage,
		Images:           b.Images,
		AuthorID:         b.AuthorID,
		Status:           b.Status,
		Slug:             b.Slug,
		LikeCount:        b.LikeCount,
		CommentCount:     b.CommentCount,
		CreatedAt:        b.CreatedAt,
		UpdatedAt:        b.UpdatedAt,
	}

	if b.Author.ID != 0 {
		res.Author = &dto.AuthorInfo{
			ID:             b.Author.ID,
			Name:           b.Author.Name,
			Email:          b.Author.Email,
			Role:           string(b.Author.Role),
			ProfileImage:   b.Author.ProfileImage,
			ProfileImageId: b.Author.ProfileImageId,
		}
	}

	return res
}
