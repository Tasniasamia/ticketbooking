package language

import (
	"ticketBooking/internal/translation"
	"ticketBooking/internal/config"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func LanguageRegisterRoutes(e *echo.Group, db *gorm.DB, config config.Config) {
	repo := NewRepository(db)
	transRepo := translation.NewRepository(db)
svc := NewService(repo, transRepo)
	h := NewHandler(svc)

	route := e.Group("/languages")
	route.POST("", h.Create)
	route.GET("", h.GetAll)
	route.GET("/:id", h.GetByID)
	route.PUT("/:id", h.Update)
	route.DELETE("/:id", h.Delete)
}