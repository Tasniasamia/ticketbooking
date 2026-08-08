package media

import (
	"ticketBooking/internal/config"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func MediaRegisterRoutes(e *echo.Group, db *gorm.DB,config config.Config) {
	cld, err := NewCloudinaryClient(config)
	if err != nil {
		// In production log & continue or panic based on preference
		panic("Cloudinary init failed: " + err.Error())
	}

	repo := NewRepository(db)
	svc := NewMediaService(repo, cld)
	h := NewHandler(svc)

	media := e.Group("/media")

	media.POST("/upload", h.UploadFile)
	media.GET("", h.List)
	media.GET("/:id", h.GetByID)
	media.DELETE("/:id", h.Delete)
}
