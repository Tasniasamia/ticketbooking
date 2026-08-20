package payment

import (
	"ticketBooking/internal/auth"
	"ticketBooking/internal/booking"
	"ticketBooking/internal/config"
	"ticketBooking/internal/currency"
	"ticketBooking/internal/event"
	middleware "ticketBooking/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func PaymentRegisterRoutes(api *echo.Group, db *gorm.DB, cfg config.Config) {
	svc := NewService(
		NewRepository(db),
		booking.NewRepository(db),
		event.NewRepository(db),
		currency.NewService(currency.NewRepository(db)),			
		
	)
	h := NewHandler(svc)
	jwtSvc := auth.NewJWTService(cfg.JwtSecret)

	g := api.Group("/payments")
	// Checkout & payment queries (auth)
	g.POST("/checkout", h.CreateCheckout, middleware.AuthMiddleware(jwtSvc))
	g.GET("", h.GetMyPayments, middleware.AuthMiddleware(jwtSvc))
    g.GET("/manager", h.GetManagementPayments, middleware.AuthMiddleware(jwtSvc),middleware.ManagerMiddleware());
	g.GET("/admin", h.GetAllPayments, middleware.AuthMiddleware(jwtSvc),middleware.AdminMiddleware());
    g.GET("/user", h.GetUserPayments, middleware.AuthMiddleware(jwtSvc));
	g.GET("/:id", h.GetPayment, middleware.AuthMiddleware(jwtSvc));
	

	g.GET("/transaction/:transaction_id", h.GetByTransactionID, middleware.AuthMiddleware(jwtSvc))
    g.GET("/verify-session", h.VerifyStripeSession);
	// register.go
    g.GET("/verify-sslcommerz", h.VerifySSLCommerzSession)
	// Webhooks (no auth — verified by signature / IPN body)
	g.POST("/webhook/stripe", h.StripeWebhook)
	g.POST("/webhook/sslcommerz", h.SSLCommerzIPN)
  




	// Payment methods CRUD
	// Public: list enabled methods (frontend payment selection)
	// Admin: create / update / delete / list all (protected)
	pm := api.Group("/payment-methods")
	pm.GET("", h.ListPaymentMethods) // ?all=true needs auth ideally; public returns enabled only
	pm.GET("/:id", h.GetPaymentMethod)
	pm.POST("", h.CreatePaymentMethod, middleware.AuthMiddleware(jwtSvc),middleware.AdminMiddleware())
	pm.PUT("/:id", h.UpdatePaymentMethod, middleware.AuthMiddleware(jwtSvc),middleware.AdminMiddleware())
	pm.PATCH("/:id", h.UpdatePaymentMethod, middleware.AuthMiddleware(jwtSvc),middleware.AdminMiddleware())
	pm.DELETE("/:id", h.DeletePaymentMethod, middleware.AuthMiddleware(jwtSvc),middleware.AdminMiddleware())
}
