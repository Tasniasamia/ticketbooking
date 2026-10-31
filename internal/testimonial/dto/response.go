package dto
import "time"
type TestimonialResponse struct {
	ID uint `json:"id"`; Name string `json:"name"`; Role string `json:"role"`; Text string `json:"text"`
	ImageURL string `json:"image_url"`; ImageID uint `json:"image_id"`; IsActive bool `json:"is_active"`
	SortOrder int `json:"sort_order"`; CreatedAt time.Time `json:"created_at"`; UpdatedAt time.Time `json:"updated_at"`
}
