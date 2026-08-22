package dto

import "time"

type UserInfo struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	ProfileImage   string `json:"profile_image,omitempty"`
	ProfileImageId uint   `json:"profile_image_id,omitempty"`
}

type Response struct {
	ID        uint        `json:"id"`
	BlogID    uint        `json:"blog_id"`
	UserID    uint        `json:"user_id"`
	User      *UserInfo   `json:"user,omitempty"`
	ParentID  *uint       `json:"parent_id"`
	Content   string      `json:"content"`
	IsDeleted bool        `json:"is_deleted"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Replies   []*Response `json:"replies"` // nested children
}
