package booking

import (
	"ticketBooking/internal/auth"
	"ticketBooking/internal/event"
	middleware "ticketBooking/internal/middlewares"
	"ticketBooking/internal/config"
    "github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	bookingRepo := NewRepository(db)
	eventRepo := event.NewRepository(db)

	svc := NewService(bookingRepo, eventRepo)
	handler := NewHandler(svc)

	jwtService := auth.NewJWTService(cfg.JwtSecret)

	api := e.Group("/api/v1/bookings", middleware.AuthMiddleware(jwtService))
	api.POST("", handler.CreateBooking)
}