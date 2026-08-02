// internal/server/routes.go
package server

import (
	"ticketBooking/internal/user"
     "github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterAllRoutes(api *echo.Group, db *gorm.DB) {
	user.UserRegisterRoutes(api, db)
	// booking.BookingRegisterRoutes(api, db)
	// ticket.TicketRegisterRoutes(api, db)
}