package messaging

import (
	"ticketBooking/internal/messaging/dto"
	"ticketBooking/internal/user"
	"time"

	"gorm.io/gorm"
)

type Conversation struct {
	gorm.Model

	User1ID uint      `json:"user1_id" gorm:"not null;index;uniqueIndex:idx_conv_pair"`
	User2ID uint      `json:"user2_id" gorm:"not null;index;uniqueIndex:idx_conv_pair"`
	User1   user.User `json:"user1,omitempty" gorm:"foreignKey:User1ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User2   user.User `json:"user2,omitempty" gorm:"foreignKey:User2ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	EventID *uint `json:"event_id" gorm:"index"`

	LastMessageAt *time.Time `json:"last_message_at"`
}

func (Conversation) TableName() string { return "conversations" }

type Message struct {
	gorm.Model

	ConversationID uint         `json:"conversation_id" gorm:"not null;index"`
	Conversation   Conversation `json:"-" gorm:"foreignKey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	SenderID uint      `json:"sender_id" gorm:"not null;index"`
	Sender   user.User `json:"sender,omitempty" gorm:"foreignKey:SenderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Content string `json:"content" gorm:"type:text"`

	AttachmentURL     string `json:"attachment_url,omitempty" gorm:"type:text"`
	AttachmentType    string `json:"attachment_type,omitempty" gorm:"type:varchar(32)"`
	AttachmentName    string `json:"attachment_name,omitempty" gorm:"type:varchar(255)"`
	AttachmentMediaID *uint  `json:"attachment_media_id,omitempty"`

	IsEdited bool `json:"is_edited" gorm:"default:false"`
	IsRead   bool `json:"is_read" gorm:"default:false;index"`
}

func (Message) TableName() string { return "messages" }

func (m *Message) ToResponse() *dto.MessageResponse {
	res := &dto.MessageResponse{
		ID:                m.ID,
		ConversationID:    m.ConversationID,
		SenderID:          m.SenderID,
		Content:           m.Content,
		IsRead:            m.IsRead,
		IsEdited:          m.IsEdited,
		AttachmentURL:     m.AttachmentURL,
		AttachmentType:    m.AttachmentType,
		AttachmentName:    m.AttachmentName,
		AttachmentMediaID: m.AttachmentMediaID,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
	if m.Sender.ID != 0 {
		res.Sender = &dto.UserBrief{
			ID:             m.Sender.ID,
			Name:           m.Sender.Name,
			Email:          m.Sender.Email,
			Role:           string(m.Sender.Role),
			ProfileImage:   m.Sender.ProfileImage,
			ProfileImageId: m.Sender.ProfileImageId,
		}
	}
	return res
}

func NormalizePair(a, b uint) (uint, uint) {
	if a < b {
		return a, b
	}
	return b, a
}
