package settings

import (
	"ticketBooking/internal/auth"
	"ticketBooking/internal/config"
	middleware "ticketBooking/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func SettingsRegisterRoutes(api *echo.Group, db *gorm.DB, cfg config.Config) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	jwtSvc := auth.NewJWTService(cfg.JwtSecret)

	g := api.Group("/settings")

	// Public – frontend needs site name, logo, social links, enabled gateways etc.
	g.GET("", h.Get)

	// Admin – create or update the single settings document
	g.PUT("", h.Upsert, middleware.AuthMiddleware(jwtSvc))
	g.POST("", h.Upsert, middleware.AuthMiddleware(jwtSvc))
}