package settings

import (
	"ticketBooking/internal/auth"
	"ticketBooking/internal/config"
	middleware "ticketBooking/internal/middlewares"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)
func SettingsRegisterRoutes(api *echo.Group, db *gorm.DB, cfg config.Config) {
	h := NewHandler(NewService(NewRepository(db)))
	jwt := auth.NewJWTService(cfg.JwtSecret)
	g := api.Group("/settings")
	g.GET("", h.Get)
	g.PUT("", h.Upsert, middleware.AuthMiddleware(jwt))
	g.POST("", h.Upsert, middleware.AuthMiddleware(jwt))
	g.GET("/pages/:slug", h.GetPageSetting)
	g.PUT("/pages/:slug", h.UpsertPageSetting, middleware.AuthMiddleware(jwt), middleware.AdminMiddleware())
}
