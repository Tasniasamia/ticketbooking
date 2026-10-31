package faq

import (
	"time"
	"gorm.io/gorm"
)

type FAQ struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Question  string         `gorm:"type:text;not null" json:"question"`
	Answer    string         `gorm:"type:text;not null" json:"answer"`
	IsActive  bool           `gorm:"default:true;not null" json:"is_active"`
	SortOrder int            `gorm:"default:0;not null" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
