package user

import (
	"net/http"
	"ticketBooking/internal/auth"
	"ticketBooking/internal/config"
	"ticketBooking/internal/email"
	"ticketBooking/internal/middlewares"
	"ticketBooking/internal/otp"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)


func UserRegisterRoutes(e *echo.Group,db *gorm.DB,config config.Config){
	NewUserRepository:=NewUserRepository(db);
	JWTService:=auth.NewJWTService(config.JwtSecret);
    emailService := email.NewEmailService(config);
     
	NewOTPRepository:=otp.NewOTPRepository(db);
	otpService:=otp.NewOTPService(NewOTPRepository,emailService)

	NewUserService:=NewUserService(NewUserRepository,JWTService,otpService);
	NewUserHandler:=NewHandler(NewUserService);

    authRoute:=e.Group("/auth");
	authRoute.POST("/register",NewUserHandler.CreateUser);
    authRoute.POST("/login",NewUserHandler.LoginUser);
    authRoute.GET("/me",NewUserHandler.GetMe,middleware.AuthMiddleware(JWTService));
	authRoute.PATCH("/me",NewUserHandler.UpdateUser,middleware.AuthMiddleware(JWTService));
	authRoute.PATCH("/members/:id", NewUserHandler.UpdateMember,middleware.AuthMiddleware(JWTService),middleware.AdminMiddleware())
	authRoute.DELETE("/:id",NewUserHandler.DeleteUser,middleware.AuthMiddleware(JWTService),middleware.AdminMiddleware());
    authRoute.POST("/verify-otp", NewUserHandler.VerifyOTP)
	authRoute.POST("/forgot-password", NewUserHandler.ForgotPassword)
	authRoute.POST("/reset-password", NewUserHandler.ResetPassword)
	authRoute.POST("/resend-otp", NewUserHandler.ResendOTP)
	authRoute.GET("/all", NewUserHandler.GetAllUsers,middleware.AuthMiddleware(JWTService),middleware.AdminMiddleware())
   	authRoute.GET("/all/managers", NewUserHandler.GetAllManager,middleware.AuthMiddleware(JWTService),middleware.AdminMiddleware())

	authRoute.GET("/:id", NewUserHandler.GetUserById,middleware.AuthMiddleware(JWTService),middleware.AdminMiddleware())
	
	authRoute.GET("/", func(c *echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
  });


}