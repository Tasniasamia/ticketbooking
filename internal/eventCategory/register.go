package eventCategory

import (
	"ticketBooking/internal/auth"
	"ticketBooking/internal/config"
	middleware "ticketBooking/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func EventCategoryRegisterRoutes(e *echo.Group, db *gorm.DB, config config.Config) {
	NewUserRepository := NewRepository(db)
    JWTService:=auth.NewJWTService(config.JwtSecret);
	NewEventCategoryService := NewEventCategoryService(NewUserRepository)
	NewHandler := NewHandler(NewEventCategoryService)

	evenCategorytRoute := e.Group("/event-categories")
	evenCategorytRoute.POST("", NewHandler.CreateEventCategory, middleware.AuthMiddleware(JWTService), middleware.ManagerMiddleware())
	evenCategorytRoute.GET("", NewHandler.GetAllEventCategories);
	evenCategorytRoute.GET("/:id", NewHandler.GetEventCategoryByID);
	evenCategorytRoute.PUT("/:id", NewHandler.UpdateEventCategory, middleware.AuthMiddleware(JWTService), middleware.ManagerMiddleware());
	evenCategorytRoute.DELETE("/:id", NewHandler.DeleteEventCategory, middleware.AuthMiddleware(JWTService), middleware.ManagerMiddleware());
}