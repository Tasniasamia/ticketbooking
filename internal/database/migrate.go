package database

import (
	"ticketBooking/internal/blog"
	"ticketBooking/internal/booking"
	"ticketBooking/internal/comment"
	"ticketBooking/internal/contact"
	"ticketBooking/internal/faq"
	"ticketBooking/internal/currency"
	"ticketBooking/internal/event"
	"ticketBooking/internal/eventCategory"
	"ticketBooking/internal/language"
	"ticketBooking/internal/media"
	"ticketBooking/internal/messaging"
	"ticketBooking/internal/otp"
	"ticketBooking/internal/payment"
	"ticketBooking/internal/settings"
	"ticketBooking/internal/subscribe"
	"ticketBooking/internal/testimonial"

	"ticketBooking/internal/translation"
	"ticketBooking/internal/user"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&user.User{},
		&event.Event{},
		&language.Language{},
		&translation.Translation{},
		&media.Media{},
		&currency.Currency{},
		&settings.Setting{},
		&settings.PageSettings{},
		&payment.Payment{},
		&booking.Booking{},
		&payment.PaymentMethod{},
		&otp.OTP{},
		&eventCategory.EventCategory{},
		&blog.Blog{},
		&blog.BlogLike{},
		&comment.Comment{},
        &messaging.Conversation{},
		&messaging.Message{},
		&faq.FAQ{},
		&testimonial.Testimonial{},
		&contact.Contact{},
		&subscribe.Subscribe{},
	)
	if err != nil {
		panic("migration failed: " + err.Error())
	}
}
