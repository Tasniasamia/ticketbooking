package server

import (
	"ticketBooking/internal/domain/event"
	"ticketBooking/internal/domain/language"
	"ticketBooking/internal/domain/translation"
	"ticketBooking/internal/domain/user"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type RouteRegistrar func(api *echo.Group, db *gorm.DB)

func RegisterAllRoutes(api *echo.Group, db *gorm.DB) {
	registrars := []RouteRegistrar{
		user.UserRegisterRoutes,
		event.EventRegisterRoutes,
		language.LanguageRegisterRoutes,
		translation.TranslationRegisterRoutes,
		// booking.RegisterRoutes,
		// ticket.RegisterRoutes,
	}

	for _, register := range registrars {
		register(api, db)
	}
}