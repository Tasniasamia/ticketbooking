package messaging

import (
	"ticketBooking/internal/auth"
	"ticketBooking/internal/config"
	middleware "ticketBooking/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(api *echo.Group, db *gorm.DB, cfg config.Config, hub *Hub) {
	repo := NewRepository(db)
	svc := NewService(repo, hub)
	jwtSvc := auth.NewJWTService(cfg.JwtSecret)
	h := NewHandler(svc, hub, jwtSvc)

	g := api.Group("/messaging")

	g.GET("/ws", h.ServeWS)

	authMW := middleware.AuthMiddleware(jwtSvc)

	g.POST("/conversations", h.StartConversation, authMW)
	g.GET("/conversations", h.GetMyConversations, authMW)
	g.GET("/conversations/:id/messages", h.GetMessages, authMW)
	g.POST("/conversations/:id/read", h.MarkRead, authMW)
	g.POST("/messages", h.SendMessage, authMW)
	g.PUT("/messages/:id", h.EditMessage, authMW)
	g.DELETE("/messages/:id", h.DeleteMessage, authMW)
}
