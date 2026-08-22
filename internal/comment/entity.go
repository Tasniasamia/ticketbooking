package comment

import (
	"ticketBooking/internal/comment/dto"
	"ticketBooking/internal/user"

	"gorm.io/gorm"
)

// Comment supports infinite nested replies via ParentID (adjacency list).
// parent_id = NULL  → top-level comment on a blog
// parent_id = X     → reply to comment X (can nest as deep as needed)
type Comment struct {
	gorm.Model

	BlogID  uint   `json:"blog_id" gorm:"not null;index"`
	UserID  uint   `json:"user_id" gorm:"not null;index"`
	User    user.User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	ParentID *uint   `json:"parent_id" gorm:"index"` // nil = root comment
	Content  string  `json:"content" gorm:"type:text;not null"`

	// optional: soft-hide without deleting children
	IsDeleted bool `json:"is_deleted" gorm:"default:false"`
}

func (c *Comment) ToResponse() *dto.Response {
	res := &dto.Response{
		ID:        c.ID,
		BlogID:    c.BlogID,
		UserID:    c.UserID,
		ParentID:  c.ParentID,
		Content:   c.Content,
		IsDeleted: c.IsDeleted,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Replies:   []*dto.Response{},
	}

	if c.User.ID != 0 {
		res.User = &dto.UserInfo{
			ID:             c.User.ID,
			Name:           c.User.Name,
			Email:          c.User.Email,
			ProfileImage:   c.User.ProfileImage,
			ProfileImageId: c.User.ProfileImageId,
		}
	}

	// if soft-deleted, hide content
	if c.IsDeleted {
		res.Content = "[deleted]"
	}

	return res
}
