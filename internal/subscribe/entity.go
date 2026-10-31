package subscribe
import ("time"; "gorm.io/gorm")
type Subscribe struct {
	ID uint `gorm:"primaryKey" json:"id"`
	Email string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	IsActive bool `gorm:"default:true;not null" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
