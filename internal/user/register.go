package user

import (
	"net/http"
    "ticketBooking/internal/auth"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
	"ticketBooking/internal/middlewares"
	"ticketBooking/internal/config"

)


func UserRegisterRoutes(e *echo.Group,db *gorm.DB,config config.Config){
	NewUserRepository:=NewUserRepository(db);
	JWTService:=auth.NewJWTService(config.JwtSecret);

	NewUserService:=NewUserService(NewUserRepository,JWTService);
	NewUserHandler:=NewHandler(NewUserService);

    authRoute:=e.Group("/auth");
	authRoute.POST("/users",NewUserHandler.CreateUser);
    authRoute.POST("/login",NewUserHandler.LoginUser);
    authRoute.GET("/me",NewUserHandler.GetMe,middleware.AuthMiddleware(JWTService));
	authRoute.PATCH("/me",NewUserHandler.UpdateUser,middleware.AuthMiddleware(JWTService));
	authRoute.DELETE("/:id",NewUserHandler.DeleteUser,middleware.AuthMiddleware(JWTService),middleware.AdminMiddleware());
authRoute.POST("/verify-otp", NewUserHandler.VerifyOTP)
	authRoute.POST("/forgot-password", NewUserHandler.ForgotPassword)
	authRoute.POST("/reset-password", NewUserHandler.ResetPassword)
	authRoute.GET("/", func(c *echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
  });


}