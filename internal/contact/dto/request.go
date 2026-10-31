package dto
type CreateContactRequest struct {
	Name string `json:"name" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
	Subject string `json:"subject" validate:"required,min=3"`
	Message string `json:"message" validate:"required,min=5"`
}
type ReplyContactRequest struct { ReplyMessage string `json:"reply_message" validate:"required,min=3"` }
