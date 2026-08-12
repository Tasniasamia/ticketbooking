package booking

import (
	"ticketBooking/internal/auth"
	"ticketBooking/internal/config"
	"ticketBooking/internal/event"
	middleware "ticketBooking/internal/middlewares"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func BookingRegisterRoutes(api *echo.Group, db *gorm.DB, cfg config.Config) {
	svc := NewService(NewRepository(db), event.NewRepository(db))
	h := NewHandler(svc)
	jwtSvc := auth.NewJWTService(cfg.JwtSecret)
	g := api.Group("/bookings", middleware.AuthMiddleware(jwtSvc))
	g.POST("", h.CreateBooking)
	g.GET("", h.GetMyBookings)
	g.GET("/:id", h.GetByID)
}
