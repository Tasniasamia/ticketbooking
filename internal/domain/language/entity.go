package language

import "gorm.io/gorm"

type Language struct {
	gorm.Model
	Name      string `json:"name" gorm:"type:varchar(50);not null"`
	Code      string `json:"code" gorm:"type:varchar(5);uniqueIndex;not null"`
	RTL       bool   `json:"rtl" gorm:"default:false"`
	IsActive  bool   `json:"is_active" gorm:"default:true"`
	IsDefault bool   `json:"is_default" gorm:"default:false"`
	Flag      string `json:"flag" gorm:"type:varchar(255)"`
}