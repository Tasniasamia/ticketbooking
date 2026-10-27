// package server;

// import (
// 	"ticketBooking/internal/config"
	
// 	"github.com/go-playground/validator/v10"
// 	"github.com/labstack/echo/v5"
// 	"github.com/labstack/echo/v5/middleware"
// 	"ticketBooking/internal/messaging"

// 	"gorm.io/gorm"
// )
// type CustomValidator struct {
//   validator *validator.Validate
// }

// func (cv *CustomValidator) Validate(i any) error {
//   if err := cv.validator.Struct(i); err != nil {
//     return echo.ErrBadRequest.Wrap(err)
//   }
//   return nil
// }

// func Start(cfg config.Config,db *gorm.DB) {


//     //server initiate
// 	e := echo.New()

// 	// message module
// 	e.Validator = &CustomValidator{validator: validator.New()}
//    	hub := messaging.NewHub()
// 	go hub.Run()
	
//   // server Build-In middleware
// 	e.Use(middleware.RequestLogger())
// 	e.Use(middleware.Recover())
//     e.Use(middleware.CORSWithConfig(
//         middleware.CORSConfig{
//             AllowOrigins: []string{
//                 "*",
//             },
//             AllowMethods: []string{
//                 "GET",
//                 "POST",
//                 "PUT",
//                 "PATCH",
//                 "DELETE",
//                 "OPTIONS",
//             },
//             AllowHeaders: []string{
//                 "Origin",
//                 "Content-Type",
//                 "Authorization",
//             },
//         },
//     ))


// 	//router middleware
// 	api:=e.Group("/api/v1");
// 	RegisterAllRoutes(api, db, cfg,hub)


   
// 	port:=cfg.Port
//   // server listening
// 	if err := e.Start(":"+port); err != nil {
// 		e.Logger.Error("failed to start server", "error", err)
// 	}
// }

package server

import (
	"context"
	"ticketBooking/internal/booking"
	"ticketBooking/internal/config"
	"ticketBooking/internal/messaging"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

func Start(cfg config.Config, db *gorm.DB) {

	//server initiate
	e := echo.New()

	// message module
	e.Validator = &CustomValidator{validator: validator.New()}
	hub := messaging.NewHub()
	go hub.Run()

	// abandoned pending booking er ticket ferot dewar background worker
	booking.StartExpiryWorker(context.Background(), booking.NewRepository(db))

	// server Build-In middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{
				"*",
			},
			AllowMethods: []string{
				"GET",
				"POST",
				"PUT",
				"PATCH",
				"DELETE",
				"OPTIONS",
			},
			AllowHeaders: []string{
				"Origin",
				"Content-Type",
				"Authorization",
			},
		},
	))

	//router middleware
	api := e.Group("/api/v1")
	RegisterAllRoutes(api, db, cfg, hub)

	port := cfg.Port
	// server listening
	if err := e.Start(":" + port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
