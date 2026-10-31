package dto
import "time"
type SubscribeResponse struct {
	ID uint `json:"id"`; Email string `json:"email"`; IsActive bool `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`; UpdatedAt time.Time `json:"updated_at"`
}
