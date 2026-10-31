package testimonial
import ("ticketBooking/internal/auth"; "ticketBooking/internal/config"; middleware "ticketBooking/internal/middlewares"; "github.com/labstack/echo/v5"; "gorm.io/gorm")
func TestimonialRegisterRoutes(api *echo.Group, db *gorm.DB, cfg config.Config) {
	h := NewHandler(NewService(NewRepository(db)))
	jwtSvc := auth.NewJWTService(cfg.JwtSecret)
	g := api.Group("/testimonials")
	g.GET("/public", h.ListPublic); g.GET("/public/:id", h.GetByID)
	g.GET("/admin", h.ListAdmin, middleware.AuthMiddleware(jwtSvc), middleware.AdminMiddleware())
	g.GET("/admin/:id", h.GetByID, middleware.AuthMiddleware(jwtSvc), middleware.AdminMiddleware())
	g.POST("", h.Create, middleware.AuthMiddleware(jwtSvc), middleware.AdminMiddleware())
	g.PUT("/:id", h.Update, middleware.AuthMiddleware(jwtSvc), middleware.AdminMiddleware())
	g.DELETE("/:id", h.Delete, middleware.AuthMiddleware(jwtSvc), middleware.AdminMiddleware())
}
