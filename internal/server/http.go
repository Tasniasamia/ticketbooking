package server;

import (
	"ticketBooking/internal/config"
	"ticketBooking/internal/user"
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

func Start(cfg config.Config,db *gorm.DB) {
	//migration 
	db.AutoMigrate(&user.User{})

    //server initiate
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
   
	
    


    // server buildin middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())



	//router middleware
	api:=e.Group("/api/v1");
	RegisterAllRoutes(api, db, cfg)


   
	port:=cfg.Port
  // server listening
	if err := e.Start(":"+port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}