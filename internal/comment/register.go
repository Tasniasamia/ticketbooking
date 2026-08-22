package comment

import (
	"ticketBooking/internal/auth"
	"ticketBooking/internal/config"
	middleware "ticketBooking/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// blogCounterFromDB — blog package import ছাড়া raw SQL দিয়ে count update
type blogCounterFromDB struct {
	db *gorm.DB
}

func (c *blogCounterFromDB) IncrementCommentCount(blogID uint) error {
	return c.db.Table("blogs").Where("id = ?", blogID).
		Update("comment_count", gorm.Expr("comment_count + 1")).Error
}

func (c *blogCounterFromDB) DecrementCommentCount(blogID uint) error {
	return c.db.Table("blogs").Where("id = ?", blogID).
		Update("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)")).Error
}

func CommentRegisterRoutes(e *echo.Group, db *gorm.DB, cfg config.Config) {
	repo := NewRepository(db)
	counter := &blogCounterFromDB{db: db}
	svc := NewCommentService(repo, counter)
	h := NewHandler(svc)
	jwtService := auth.NewJWTService(cfg.JwtSecret)

	commentRoute := e.Group("/comments")

	commentRoute.GET("/blog/:blogId", h.GetCommentsByBlog)
	commentRoute.POST("", h.CreateComment, middleware.AuthMiddleware(jwtService))
	commentRoute.PUT("/:id", h.UpdateComment, middleware.AuthMiddleware(jwtService))
	commentRoute.DELETE("/:id", h.DeleteComment, middleware.AuthMiddleware(jwtService))
}
