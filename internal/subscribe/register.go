package subscribe
import ("ticketBooking/internal/auth"; "ticketBooking/internal/config"; middleware "ticketBooking/internal/middlewares"; "github.com/labstack/echo/v5"; "gorm.io/gorm")
func SubscribeRegisterRoutes(api *echo.Group, db *gorm.DB, cfg config.Config) {
	h := NewHandler(NewService(NewRepository(db)))
	jwtSvc := auth.NewJWTService(cfg.JwtSecret)
	g := api.Group("/subscribers")
	g.POST("", h.Create)
	g.GET("/admin", h.ListAdmin, middleware.AuthMiddleware(jwtSvc), middleware.AdminMiddleware())
}
