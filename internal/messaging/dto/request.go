package dto

type SendMessageRequest struct {
	ReceiverID uint   `json:"receiver_id" validate:"required"`
	Content    string `json:"content" validate:"required,min=1,max=5000"`
	EventID    *uint  `json:"event_id,omitempty"` // optional context
}

type StartConversationRequest struct {
	ReceiverID uint  `json:"receiver_id" validate:"required"`
	EventID    *uint `json:"event_id,omitempty"`
	UserId     uint  `json:"user_id" `
}

type MarkReadRequest struct {
	MessageIDs []uint `json:"message_ids" validate:"required,min=1"`
}
