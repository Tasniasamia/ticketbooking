// package event

// import (
// 	"ticketBooking/internal/auth"
// 	"ticketBooking/internal/config"
// 	middleware "ticketBooking/internal/middlewares"
// 	"ticketBooking/internal/user"

// 	"github.com/labstack/echo/v5"
// 	"gorm.io/gorm"
// )

// func EventRegisterRoutes(e *echo.Group, db *gorm.DB, config config.Config) {
// 	NewUserRepository := NewRepository(db);
// 	UserRepository :=user.NewUserRepository(db);
//     JWTService:=auth.NewJWTService(config.JwtSecret);
 
// 	NewEventService := NewEventService(NewUserRepository, UserRepository)
// 	NewHandler := NewHandler(NewEventService)

// 	eventRoute := e.Group("/events")
// 	eventRoute.POST("", NewHandler.CreateEvent,middleware.AuthMiddleware(JWTService),middleware.ManagerMiddleware());
// 	eventRoute.GET("", NewHandler.GetAllEvents);
// 	eventRoute.GET("/:id", NewHandler.GetEventByID);
// 	eventRoute.GET("/admin/:id", NewHandler.GetEventByIDAdmin,middleware.AuthMiddleware(JWTService),middleware.AdminMiddleware());
// 	eventRoute.GET("/myEvents", NewHandler.GetMyEvents,middleware.AuthMiddleware(JWTService));
// 	eventRoute.PUT("/:id", NewHandler.UpdateEvent,middleware.AuthMiddleware(JWTService),middleware.ManagerMiddleware());
// 	eventRoute.DELETE("/:id", NewHandler.DeleteEvent,middleware.AuthMiddleware(JWTService),middleware.ManagerMiddleware());
//     eventRoute.PATCH("/updateStatus", NewHandler.UpdateEventStatus,middleware.AuthMiddleware(JWTService),middleware.AdminMiddleware());
// 	eventRoute.GET("/public", NewHandler.GetAllEventsPublic);

// }

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
	repo := NewRepository(db)
	userRepo := user.NewUserRepository(db)
	jwtService := auth.NewJWTService(config.JwtSecret)

	service := NewEventService(repo, userRepo)
	h := NewHandler(service)

	authMW := middleware.AuthMiddleware(jwtService)
	manager := middleware.ManagerMiddleware()
	admin := middleware.AdminMiddleware()
	adminOrManager := middleware.AdminOrManagerMiddleware()

	route := e.Group("/events")

	// ---- Public: শুধু approved event ----
	route.GET("", h.GetAllEventsPublic)
	route.GET("/public", h.GetAllEventsPublic)
	route.GET("/:id", h.GetEventByID)

	// ---- Manager: নিজের event ----
	route.POST("", h.CreateEvent, authMW, manager)
	route.GET("/myEvents", h.GetMyEvents, authMW, manager)

	// ---- Admin: সব event, যেকোনো status ----
	route.GET("/admin", h.GetAllEventsAdmin, authMW, admin)
	route.GET("/admin/:id", h.GetEventByIDAdmin, authMW, admin)
	route.PATCH("/updateStatus", h.UpdateEventStatus, authMW, admin)

	// ---- Update / Delete: manager (শুধু নিজের) অথবা admin ----
	route.PUT("/:id", h.UpdateEvent, authMW, adminOrManager)
	route.DELETE("/:id", h.DeleteEvent, authMW, adminOrManager)
}
