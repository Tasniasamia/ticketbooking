package event

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
	"ticketBooking/internal/config"
)

func EventRegisterRoutes(e *echo.Group, db *gorm.DB, config config.Config) {
	NewUserRepository := NewRepository(db)

	NewEventService := NewEventService(NewUserRepository)
	NewHandler := NewHandler(NewEventService)

	eventRoute := e.Group("/events")
	eventRoute.POST("", NewHandler.CreateEvent)
	eventRoute.GET("", NewHandler.GetAllEvents)
}