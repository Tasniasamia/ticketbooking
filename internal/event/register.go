package event

import (
	"ticketBooking/internal/auth"
	"ticketBooking/internal/config"
	middleware "ticketBooking/internal/middlewares"
	"ticketBooking/internal/user"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func EventRegisterRoutes(e *echo.Group, db *gorm.DB, config config.Config) {
	NewUserRepository := NewRepository(db);
	UserRepository :=user.NewUserRepository(db);
    JWTService:=auth.NewJWTService(config.JwtSecret);
 
	NewEventService := NewEventService(NewUserRepository, UserRepository)
	NewHandler := NewHandler(NewEventService)

	eventRoute := e.Group("/events")
	eventRoute.POST("", NewHandler.CreateEvent,middleware.AuthMiddleware(JWTService),middleware.ManagerMiddleware());
	eventRoute.GET("", NewHandler.GetAllEvents);
	eventRoute.GET("/:id", NewHandler.GetEventByID);
	eventRoute.GET("/admin/:id", NewHandler.GetEventByIDAdmin,middleware.AuthMiddleware(JWTService),middleware.AdminMiddleware());
	eventRoute.GET("/myEvents", NewHandler.GetMyEvents,middleware.AuthMiddleware(JWTService));
	eventRoute.PUT("/:id", NewHandler.UpdateEvent,middleware.AuthMiddleware(JWTService),middleware.ManagerMiddleware());
	eventRoute.DELETE("/:id", NewHandler.DeleteEvent,middleware.AuthMiddleware(JWTService),middleware.ManagerMiddleware());
    eventRoute.PATCH("/updateStatus", NewHandler.UpdateEventStatus,middleware.AuthMiddleware(JWTService),middleware.AdminMiddleware());
	eventRoute.GET("/public", NewHandler.GetAllEventsPublic);

}