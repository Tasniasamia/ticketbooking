package currency

import (
	"ticketBooking/internal/auth"
	"ticketBooking/internal/config"
	middleware "ticketBooking/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func CurrencyRegisterRoutes(api *echo.Group, db *gorm.DB, cfg config.Config) {
	repo := NewRepository(db)
	svc := NewService(repo)
	_ = svc.SeedDefaults()
	h := NewHandler(svc)

	jwtSvc := auth.NewJWTService(cfg.JwtSecret)

	g := api.Group("/currencies")

	// Public
	g.GET("", h.List)
	g.GET("/default", h.GetDefault)
	g.GET("/:code", h.GetByCode)
	g.POST("/convert", h.Convert)

	// Admin (auth required — later add admin role middleware if needed)
	g.POST("", h.Create, middleware.AuthMiddleware(jwtSvc))
	g.PUT("/:id", h.Update, middleware.AuthMiddleware(jwtSvc))
	g.POST("/set-default", h.SetDefault, middleware.AuthMiddleware(jwtSvc))
	g.DELETE("/:id", h.Delete, middleware.AuthMiddleware(jwtSvc))
}
