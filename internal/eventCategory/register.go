// package eventCategory

// import (
// 	"ticketBooking/internal/auth"
// 	"ticketBooking/internal/config"
// 	middleware "ticketBooking/internal/middlewares"

// 	"github.com/labstack/echo/v5"
// 	"gorm.io/gorm"
// )

// func EventCategoryRegisterRoutes(e *echo.Group, db *gorm.DB, config config.Config) {
// 	NewUserRepository := NewRepository(db)
//     JWTService:=auth.NewJWTService(config.JwtSecret);
// 	NewEventCategoryService := NewEventCategoryService(NewUserRepository)
// 	NewHandler := NewHandler(NewEventCategoryService)

// 	evenCategorytRoute := e.Group("/event-categories")
// 	evenCategorytRoute.POST("", NewHandler.CreateEventCategory, middleware.AuthMiddleware(JWTService))
// 	evenCategorytRoute.GET("", NewHandler.GetAllEventCategories);
// 	evenCategorytRoute.GET("/:id", NewHandler.GetEventCategoryByID);
// 	evenCategorytRoute.PUT("/:id", NewHandler.UpdateEventCategory, middleware.AuthMiddleware(JWTService));
// 	evenCategorytRoute.DELETE("/:id", NewHandler.DeleteEventCategory, middleware.AuthMiddleware(JWTService));
// }


package eventCategory

import (
	"ticketBooking/internal/auth"
	"ticketBooking/internal/config"
	middleware "ticketBooking/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func EventCategoryRegisterRoutes(e *echo.Group, db *gorm.DB, config config.Config) {
	repo := NewRepository(db)
	jwtService := auth.NewJWTService(config.JwtSecret)
	service := NewEventCategoryService(repo)
	h := NewHandler(service)

	auth := middleware.AuthMiddleware(jwtService)
	adminOrManager := middleware.AdminOrManagerMiddleware()

	route := e.Group("/event-categories")
	route.GET("", h.GetAllEventCategories)
	route.GET("/:id", h.GetEventCategoryByID)
	route.POST("", h.CreateEventCategory, auth, adminOrManager)
	route.PUT("/:id", h.UpdateEventCategory, auth, adminOrManager)
	route.DELETE("/:id", h.DeleteEventCategory, auth, adminOrManager)
}