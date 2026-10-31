package dto
type CreateSubscribeRequest struct { Email string `json:"email" validate:"required,email"` }
