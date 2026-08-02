package config

import (
	"ticketBooking/internal/user"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(
		&user.User{},
		// অন্য models...
	)
}