package contact
import ("time"; "gorm.io/gorm")
type ContactStatus string
const (StatusPending ContactStatus = "pending"; StatusReplied ContactStatus = "replied")
type Contact struct {
	ID uint `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(150);not null" json:"name"`
	Email string `gorm:"type:varchar(255);not null;index" json:"email"`
	Subject string `gorm:"type:varchar(255);not null" json:"subject"`
	Message string `gorm:"type:text;not null" json:"message"`
	Status ContactStatus `gorm:"type:varchar(30);not null;default:'pending'" json:"status"`
	ReplyMessage string `gorm:"type:text" json:"reply_message"`
	RepliedAt *time.Time `json:"replied_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
