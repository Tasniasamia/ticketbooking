package dto
import "time"
type ContactResponse struct {
	ID uint `json:"id"`; Name string `json:"name"`; Email string `json:"email"`
	Subject string `json:"subject"`; Message string `json:"message"`; Status string `json:"status"`
	ReplyMessage string `json:"reply_message,omitempty"`; RepliedAt *time.Time `json:"replied_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`; UpdatedAt time.Time `json:"updated_at"`
}
