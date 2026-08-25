package messaging

import (
	"ticketBooking/internal/auth"
	"ticketBooking/internal/config"
	middleware "ticketBooking/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// RegisterRoutes sets up REST + WebSocket routes for messaging.
// Call this from server.RegisterAllRoutes and also start the Hub.
func RegisterRoutes(api *echo.Group, db *gorm.DB, cfg config.Config, hub *Hub) {
	repo := NewRepository(db)
	svc := NewService(repo, hub)
	jwtSvc := auth.NewJWTService(cfg.JwtSecret)
	h := NewHandler(svc, hub, jwtSvc)

	g := api.Group("/messaging")

	// WebSocket – token via query ?token=... or Authorization header
	g.GET("/ws", h.ServeWS)

	// REST (all require auth)
	authMW := middleware.AuthMiddleware(jwtSvc)

	g.POST("/conversations", h.StartConversation)
	g.GET("/conversations", h.GetMyConversations, authMW)
	g.GET("/conversations/:id/messages", h.GetMessages, authMW)
	g.POST("/conversations/:id/read", h.MarkRead, authMW)
	g.POST("/messages", h.SendMessage, authMW)
}
