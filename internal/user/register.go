package user

import (
	"net/http"
    "ticketBooking/internal/auth"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)


func UserRegisterRoutes(e *echo.Group,db *gorm.DB){
	NewUserRepository:=NewUserRepository(db);
	JWTService:=auth.NewJWTService("");

	NewUserService:=NewUserService(NewUserRepository,JWTService);
	NewUserHandler:=NewHandler(NewUserService);

    authRoute:=e.Group("/auth");
	authRoute.POST("/users",NewUserHandler.CreateUser);
    authRoute.POST("/login",NewUserHandler.LoginUser);


	authRoute.GET("/", func(c *echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
  });


}