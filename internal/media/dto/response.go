package dto

type Response struct {
	ID           uint   `json:"id"` // auto-generated DB id  → use this to find & delete later
	ImageID      string `json:"image_id"`                   // Cloudinary public_id → used to destroy on Cloudinary
	PublicURL    string `json:"public_url"`
	SecureURL    string `json:"secure_url"`
	ResourceType string `json:"resource_type"`
	Format       string `json:"format"`
	Type         string `json:"type"` // image | video | pdf | file
	ModelName    string `json:"model_name"`
	ModelID      *uint  `json:"model_id,omitempty"`
	Folder       string `json:"folder,omitempty"`
	Size         int64  `json:"size"`
	Width        int    `json:"width,omitempty"`
	Height       int    `json:"height,omitempty"`
	OriginalName string `json:"original_name,omitempty"`
	CreatedAt    string `json:"created_at"`
}

// Convenience response the user asked for
// Example when uploading for event:
//
//	{
//	  "event_url": "https://res.cloudinary.com/..../image.jpg",
//	  "event_id":  15,          // media table auto id
//	  "image_id":  "events/abc",
//	  "type":      "image"
//	}
type ModelKeyedResponse struct {
	EventURL  string `json:"event_url,omitempty"`
	EventID   uint   `json:"event_id,omitempty"`
	UserURL   string `json:"user_url,omitempty"`
	UserID    uint   `json:"user_id,omitempty"`
	ProductURL string `json:"product_url,omitempty"`
	ProductID uint   `json:"product_id,omitempty"`
	ImageID   string `json:"image_id"`
	PublicURL string `json:"public_url"`
	Type      string `json:"type"`
	ModelName string `json:"model_name"`
	CreatedAt string `json:"created_at"`
}
