// package user

// import (
// 	"net/http"
// 	"ticketBooking/internal/auth"
// 	"ticketBooking/internal/config"
// 	"ticketBooking/internal/email"
// 	"ticketBooking/internal/messaging"
// 	messageDTO "ticketBooking/internal/messaging/dto"
// 	"ticketBooking/internal/otp"
// middlewares "ticketBooking/internal/middlewares"
// 	"github.com/labstack/echo/v5"
// 	"gorm.io/gorm"
// )

// // messagingConvBridge — messaging.Service কে user.conversationStarter হিসেবে adapt করে
// // (user/service এ messaging import করলে cycle হয়; এখানে wiring OK)
// type messagingConvBridge struct {
// 	svc *messaging.Service
// }

// func (b *messagingConvBridge) Start(senderID, receiverID uint) error {
// 	_, err := b.svc.StartConversation(senderID, &messageDTO.StartConversationRequest{
// 		ReceiverID: receiverID,
// 	})
// 	return err
// }

// func UserRegisterRoutes(e *echo.Group, db *gorm.DB, config config.Config) {
// 	NewUserRepository := NewUserRepository(db)
// 	JWTService := auth.NewJWTService(config.JwtSecret)
// 	emailService := email.NewEmailService(config)

// 	NewOTPRepository := otp.NewOTPRepository(db)
// 	otpService := otp.NewOTPService(NewOTPRepository, emailService)

// 	// Hub ছাড়া service — শুধু conversation create (WS লাগে না verify flow-এ)
// 	msgSvc := messaging.NewService(messaging.NewRepository(db), nil)
// 	convBridge := &messagingConvBridge{svc: msgSvc}

// 	NewUserService := NewUserService(NewUserRepository, JWTService, otpService, convBridge)
// 	NewUserHandler := NewHandler(NewUserService)

// 	authRoute := e.Group("/auth")
// 	authRoute.POST("/register", NewUserHandler.CreateUser)
// 	authRoute.POST("/login", NewUserHandler.LoginUser)
// 	authRoute.GET("/me", NewUserHandler.GetMe, middlewares.AuthMiddleware(JWTService))
// 	authRoute.PATCH("/me", NewUserHandler.UpdateUser, middlewares.AuthMiddleware(JWTService))
// 	authRoute.PATCH("/members/:id", NewUserHandler.UpdateMember, middlewares.AuthMiddleware(JWTService), middlewares.AdminMiddleware())
// 	authRoute.DELETE("/:id", NewUserHandler.DeleteUser, middlewares.AuthMiddleware(JWTService), middlewares.AdminMiddleware())
// 	authRoute.POST("/verify-otp", NewUserHandler.VerifyOTP)
// 	authRoute.POST("/forgot-password", NewUserHandler.ForgotPassword)
// 	authRoute.POST("/reset-password", NewUserHandler.ResetPassword)
// 	authRoute.POST("/resend-otp", NewUserHandler.ResendOTP)
// 	authRoute.GET("/all", NewUserHandler.GetAllUsers, middlewares.AuthMiddleware(JWTService), middlewares.AdminMiddleware())
// 	authRoute.GET("/all/managers", NewUserHandler.GetAllManager, middlewares.AuthMiddleware(JWTService), middlewares.AdminMiddleware())
// 	authRoute.GET("/:id", NewUserHandler.GetUserById, middlewares.AuthMiddleware(JWTService), middlewares.AdminMiddleware())

// 	authRoute.GET("/", func(c *echo.Context) error {
// 		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
// 	})
// }


package user

import (
	"net/http"
	"ticketBooking/internal/auth"
	"ticketBooking/internal/config"
	"ticketBooking/internal/email"
	middleware "ticketBooking/internal/middlewares"
	"ticketBooking/internal/otp"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func UserRegisterRoutes(e *echo.Group, db *gorm.DB, config config.Config, conv ...ConversationStarter) {
	NewUserRepository := NewUserRepository(db)
	JWTService := auth.NewJWTService(config.JwtSecret)
	emailService := email.NewEmailService(config)

	NewOTPRepository := otp.NewOTPRepository(db)
	otpService := otp.NewOTPService(NewOTPRepository, emailService)

	var starter ConversationStarter
	if len(conv) > 0 {
		starter = conv[0]
	}

	NewUserService := NewUserService(NewUserRepository, JWTService, otpService, starter)
	NewUserHandler := NewHandler(NewUserService)

	authRoute := e.Group("/auth")
	authRoute.POST("/register", NewUserHandler.CreateUser)
	authRoute.POST("/login", NewUserHandler.LoginUser)
	authRoute.GET("/me", NewUserHandler.GetMe, middleware.AuthMiddleware(JWTService))
	authRoute.PATCH("/me", NewUserHandler.UpdateUser, middleware.AuthMiddleware(JWTService))
	authRoute.PATCH("/members/:id", NewUserHandler.UpdateMember, middleware.AuthMiddleware(JWTService), middleware.AdminMiddleware())
	authRoute.DELETE("/:id", NewUserHandler.DeleteUser, middleware.AuthMiddleware(JWTService), middleware.AdminMiddleware())
	authRoute.POST("/verify-otp", NewUserHandler.VerifyOTP)
	authRoute.POST("/forgot-password", NewUserHandler.ForgotPassword)
	authRoute.POST("/reset-password", NewUserHandler.ResetPassword)
	authRoute.POST("/resend-otp", NewUserHandler.ResendOTP)
	authRoute.GET("/all", NewUserHandler.GetAllUsers, middleware.AuthMiddleware(JWTService), middleware.AdminMiddleware())
	authRoute.GET("/all/managers", NewUserHandler.GetAllManager, middleware.AuthMiddleware(JWTService), middleware.AdminMiddleware())
	authRoute.GET("/:id", NewUserHandler.GetUserById, middleware.AuthMiddleware(JWTService), middleware.AdminMiddleware())

	authRoute.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
	})
}