package testimonial
import ("time"; "gorm.io/gorm")
type Testimonial struct {
	ID uint `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(150);not null" json:"name"`
	Role string `gorm:"type:varchar(150)" json:"role"`
	Text string `gorm:"type:text;not null" json:"text"`
	ImageURL string `gorm:"type:varchar(500)" json:"image_url"`
	ImageID uint `gorm:"default:0" json:"image_id"`
	IsActive bool `gorm:"default:true;not null" json:"is_active"`
	SortOrder int `gorm:"default:0;not null" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
