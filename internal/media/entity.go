package media

import (
	"ticketBooking/internal/media/dto"
	"time"

	"gorm.io/gorm"
)

// Media holds every uploaded file (image / video / pdf / raw) in one place.
// image_id  = Cloudinary public_id  → used for delete on Cloudinary
// public_url = the URL that clients actually use
type Media struct {
	gorm.Model // ID, CreatedAt, UpdatedAt, DeletedAt (soft delete)

	ImageID      string `json:"image_id" gorm:"column:image_id;type:varchar(512);not null;index"` // Cloudinary public_id
	PublicURL    string `json:"public_url" gorm:"type:text;not null"`
	SecureURL    string `json:"secure_url" gorm:"type:text"`
	ResourceType string `json:"resource_type" gorm:"type:varchar(32);not null"` // image | video | raw
	Format       string `json:"format" gorm:"type:varchar(32)"`                 // jpg, png, pdf, mp4 ...
	Type         string `json:"type" gorm:"type:varchar(32);not null;index"`    // image | video | pdf | file
	ModelName    string `json:"model_name" gorm:"type:varchar(64);not null;index"` // event | user | product | ...
	ModelID      *uint  `json:"model_id" gorm:"index"`                            // optional related entity id
	Folder       string `json:"folder" gorm:"type:varchar(255)"`
	Size         int64  `json:"size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	OriginalName string `json:"original_name" gorm:"type:varchar(512)"`
}

func (m *Media) ToResponse() *dto.Response {
	return &dto.Response{
		ID:           m.ID,
		ImageID:      m.ImageID,
		PublicURL:    m.PublicURL,
		SecureURL:    m.SecureURL,
		ResourceType: m.ResourceType,
		Format:       m.Format,
		Type:         m.Type,
		ModelName:    m.ModelName,
		ModelID:      m.ModelID,
		Folder:       m.Folder,
		Size:         m.Size,
		Width:        m.Width,
		Height:       m.Height,
		OriginalName: m.OriginalName,
		CreatedAt:    m.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// Helper for soft-delete awareness
func (m *Media) IsDeleted() bool {
	return m.DeletedAt.Valid
}

// TableName explicit
func (Media) TableName() string {
	return "media"
}

// optional convenience for response that frontend asked for (event_URL style)
func (m *Media) ToModelKeyedResponse() map[string]interface{} {
	keyURL := m.ModelName + "_url"
	keyID := m.ModelName + "_id"

	return map[string]interface{}{
		keyURL:       m.SecureURL,
		keyID:        m.ID,
		"image_id":   m.ImageID,
		"public_url": m.PublicURL,
		"type":       m.Type,
		"model_name": m.ModelName,
		"created_at": m.CreatedAt.Format(time.RFC3339),
	}
}
