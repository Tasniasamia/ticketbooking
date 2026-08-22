package server

import (
	"ticketBooking/internal/blog"
	"ticketBooking/internal/booking"
	"ticketBooking/internal/comment"

	"ticketBooking/internal/config"
	"ticketBooking/internal/currency"
	"ticketBooking/internal/event"
	"ticketBooking/internal/eventCategory"
	"ticketBooking/internal/language"
	"ticketBooking/internal/media"
	"ticketBooking/internal/payment"
	"ticketBooking/internal/settings"
	"ticketBooking/internal/translation"
	"ticketBooking/internal/user"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type RouteRegistrar func(api *echo.Group, db *gorm.DB, config config.Config)

func RegisterAllRoutes(api *echo.Group, db *gorm.DB, config config.Config) {
	registrars := []RouteRegistrar{
		user.UserRegisterRoutes,
		event.EventRegisterRoutes,
		language.LanguageRegisterRoutes,
		translation.TranslationRegisterRoutes,
		media.MediaRegisterRoutes,
		currency.CurrencyRegisterRoutes,
		settings.SettingsRegisterRoutes,
		payment.PaymentRegisterRoutes,
		booking.BookingRegisterRoutes,
		eventCategory.EventCategoryRegisterRoutes,
		blog.BlogRegisterRoutes,
		comment.CommentRegisterRoutes,
	}

	for _, register := range registrars {
		register(api, db, config)
	}
}
