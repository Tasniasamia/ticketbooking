package dashboard

import (
	"ticketBooking/internal/auth"
	"ticketBooking/internal/config"
	middleware "ticketBooking/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// routes.go এর registrars list এ আগে থেকেই আছে: dashboard.DashboardRegisterRoutes
//
//	GET /dashboard/admin    (admin)
//	GET /dashboard/manager  (manager)
//	GET /dashboard/me       (login করা যেকোনো user)
func DashboardRegisterRoutes(e *echo.Group, db *gorm.DB, config config.Config) {
	svc := NewService(NewRepository(db))
	h := NewHandler(svc)
	jwtService := auth.NewJWTService(config.JwtSecret)

	authMW := middleware.AuthMiddleware(jwtService)

	d := e.Group("/dashboard")
	d.GET("/admin", h.Admin, authMW, middleware.AdminMiddleware())
	d.GET("/manager", h.Manager, authMW, middleware.ManagerMiddleware())
	d.GET("/me", h.User, authMW)
}
