package contact
import ("ticketBooking/internal/auth"; "ticketBooking/internal/config"; "ticketBooking/internal/email"; middleware "ticketBooking/internal/middlewares"; "github.com/labstack/echo/v5"; "gorm.io/gorm")
func ContactRegisterRoutes(api *echo.Group, db *gorm.DB, cfg config.Config) {
	h := NewHandler(NewService(NewRepository(db), email.NewEmailService(cfg)))
	jwtSvc := auth.NewJWTService(cfg.JwtSecret)
	g := api.Group("/contacts")
	g.POST("", h.Create)
	g.GET("/admin", h.ListAdmin, middleware.AuthMiddleware(jwtSvc), middleware.AdminMiddleware())
	g.GET("/admin/:id", h.GetByID, middleware.AuthMiddleware(jwtSvc), middleware.AdminMiddleware())
	g.POST("/admin/:id/reply", h.Reply, middleware.AuthMiddleware(jwtSvc), middleware.AdminMiddleware())
}
