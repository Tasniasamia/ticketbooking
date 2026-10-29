package dto

type SendMessageRequest struct {
	ReceiverID        uint   `json:"receiver_id" validate:"required"`
	Content           string `json:"content" validate:"omitempty,max=5000"`
	EventID           *uint  `json:"event_id,omitempty"`
	AttachmentURL     string `json:"attachment_url,omitempty"`
	AttachmentType    string `json:"attachment_type,omitempty"`
	AttachmentName    string `json:"attachment_name,omitempty"`
	AttachmentMediaID *uint  `json:"attachment_media_id,omitempty"`
}

type StartConversationRequest struct {
	ReceiverID uint  `json:"receiver_id" validate:"required"`
	EventID    *uint `json:"event_id,omitempty"`
	UserId     uint  `json:"user_id"`
}

type MarkReadRequest struct {
	MessageIDs []uint `json:"message_ids" validate:"required,min=1"`
}

type EditMessageRequest struct {
	Content string `json:"content" validate:"required,min=1,max=5000"`
}
