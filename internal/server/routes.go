package server

import (
	"ticketBooking/internal/event"
	"ticketBooking/internal/language"
	"ticketBooking/internal/media"
	"ticketBooking/internal/translation"
	"ticketBooking/internal/user"
   "ticketBooking/internal/config"
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
		// booking.RegisterRoutes,
		// ticket.RegisterRoutes,
	}

	for _, register := range registrars {
		register(api, db, config)
	}
}
