package config

import (
	"ticketBooking/internal/event"
	"ticketBooking/internal/language"
	"ticketBooking/internal/translation"
	"ticketBooking/internal/user"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(
		&user.User{},
		&event.Event{},
		&language.Language{},
		&translation.Translation{},
	)
}