package dto

import "time"

type UserBrief struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	ProfileImage   string `json:"profile_image,omitempty"`
	ProfileImageId uint   `json:"profile_image_id,omitempty"`
}

type MessageResponse struct {
	ID                uint       `json:"id"`
	ConversationID    uint       `json:"conversation_id"`
	SenderID          uint       `json:"sender_id"`
	Sender            *UserBrief `json:"sender,omitempty"`
	Content           string     `json:"content"`
	IsRead            bool       `json:"is_read"`
	IsEdited          bool       `json:"is_edited"`
	AttachmentURL     string     `json:"attachment_url,omitempty"`
	AttachmentType    string     `json:"attachment_type,omitempty"`
	AttachmentName    string     `json:"attachment_name,omitempty"`
	AttachmentMediaID *uint      `json:"attachment_media_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type ConversationResponse struct {
	ID          uint             `json:"id"`
	Participant *UserBrief       `json:"participant"`
	EventID     *uint            `json:"event_id,omitempty"`
	LastMessage *MessageResponse `json:"last_message,omitempty"`
	UnreadCount int64            `json:"unread_count"`
	UpdatedAt   time.Time        `json:"updated_at"`
	CreatedAt   time.Time        `json:"created_at"`
}

type WSIncoming struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type WSOutgoing struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type WSMessagePayload struct {
	ConversationID uint   `json:"conversation_id"`
	ReceiverID     uint   `json:"receiver_id"`
	Content        string `json:"content"`
	EventID        *uint  `json:"event_id,omitempty"`
}

type WSReadPayload struct {
	ConversationID uint   `json:"conversation_id"`
	MessageIDs     []uint `json:"message_ids"`
}
