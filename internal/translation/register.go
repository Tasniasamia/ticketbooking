package translation

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func TranslationRegisterRoutes(e *echo.Group, db *gorm.DB) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	route := e.Group("/translations")

	route.POST("", h.Create)
	route.GET("", h.GetAll)
	route.GET("/group", h.GetGroup)
	route.PUT("/bulk", h.BulkUpdate)
	route.GET("/:key", h.GetByKey)
	route.PUT("/:key", h.UpdateByKey)
	route.DELETE("/:key", h.DeleteByKey)
}