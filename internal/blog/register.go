package blog

import (
	"ticketBooking/internal/auth"
	"ticketBooking/internal/comment"
	"ticketBooking/internal/config"
	middleware "ticketBooking/internal/middlewares"
	"ticketBooking/internal/user"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type commentTreeAdapter struct {
	svc *comment.Service
}

func (a *commentTreeAdapter) GetCommentsTree(blogID uint) (interface{}, error) {
	return a.svc.GetCommentsTree(blogID)
}

func BlogRegisterRoutes(e *echo.Group, db *gorm.DB, cfg config.Config) {
	blogRepo := NewRepository(db)
	userRepo := user.NewUserRepository(db)
	jwtService := auth.NewJWTService(cfg.JwtSecret)

	blogService := NewBlogService(blogRepo, userRepo)

	// comments tree + comment_count update
	commentRepo := comment.NewRepository(db)
	commentSvc := comment.NewCommentService(commentRepo, blogRepo)
	blogService.SetCommentLoader(&commentTreeAdapter{svc: commentSvc})

	h := NewHandler(blogService)

	blogRoute := e.Group("/blogs")

	blogRoute.GET("/public", h.GetAllBlogsPublic)
	blogRoute.GET("/:id", h.GetBlogByID)

	blogRoute.POST("", h.CreateBlog, middleware.AuthMiddleware(jwtService))

	blogRoute.GET("/admin", h.GetAllBlogs, middleware.AuthMiddleware(jwtService), middleware.AdminMiddleware())

	blogRoute.GET("/admin/:id", h.GetBlogByIDAdmin, middleware.AuthMiddleware(jwtService), middleware.AdminMiddleware())
	blogRoute.GET("/myBlogs", h.GetMyBlogs, middleware.AuthMiddleware(jwtService))
	blogRoute.PUT("/:id", h.UpdateBlog, middleware.AuthMiddleware(jwtService))

	blogRoute.DELETE("/:id", h.DeleteBlog, middleware.AuthMiddleware(jwtService), middleware.AdminMiddleware())
	blogRoute.PATCH("/admin/updateStatus", h.UpdateBlogStatus, middleware.AuthMiddleware(jwtService), middleware.AdminMiddleware())
	blogRoute.POST("/:id/like", h.ToggleLike, middleware.AuthMiddleware(jwtService))
}
