package database

import (
	"ticketBooking/internal/booking"
	"ticketBooking/internal/currency"
	"ticketBooking/internal/event"
	"ticketBooking/internal/eventCategory"
	"ticketBooking/internal/language"
	"ticketBooking/internal/media"
	"ticketBooking/internal/otp"
	"ticketBooking/internal/payment"
	"ticketBooking/internal/settings"

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
		&settings.Setting{},
		&payment.Payment{},	
		&booking.Booking{},
        &payment.PaymentMethod{},
		&otp.OTP{},
		&eventCategory.EventCategory{},
	)
}
