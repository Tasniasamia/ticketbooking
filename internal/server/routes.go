package server

import (
	"ticketBooking/internal/event"
	"ticketBooking/internal/user"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type RouteRegistrar func(api *echo.Group, db *gorm.DB)

func RegisterAllRoutes(api *echo.Group, db *gorm.DB) {
	registrars := []RouteRegistrar{
		user.UserRegisterRoutes,
		event.EventRegisterRoutes,
		// booking.RegisterRoutes,
		// ticket.RegisterRoutes,
	}

	for _, register := range registrars {
		register(api, db)
	}
}