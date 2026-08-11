package database

import (
	"ticketBooking/internal/currency"
	"ticketBooking/internal/event"
	"ticketBooking/internal/language"
	"ticketBooking/internal/media"
	
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
		&media.Media{},
		&currency.Currency{},

	)
}
