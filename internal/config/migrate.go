package config

import (
	"ticketBooking/internal/domain/event"
	"ticketBooking/internal/domain/language"
	"ticketBooking/internal/domain/translation"
	"ticketBooking/internal/domain/user"

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