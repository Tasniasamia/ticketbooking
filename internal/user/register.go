package user

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)


func UserRegisterRoutes(e *echo.Echo,db *gorm.DB){
	NewUserRepository:=NewUserRepository(db);
	NewUserService:=NewUserService(NewUserRepository);
	NewUserHandler:=NewHandler(NewUserService);

	e.POST("/users",NewUserHandler.CreateUser);
		  e.GET("/", func(c *echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
  });

}