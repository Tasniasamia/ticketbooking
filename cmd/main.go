package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


type User struct {
	gorm.Model
	Name     string `json:"name" validate:"required" gorm:"type:varchar(100);not null"`
	Email    string `json:"email" validate:"required,email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password string `string" json:"password" validate:"required,min=6" gorm:"type:varchar(100);not null"`
}
type CustomValidator struct {
  validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
  if err := cv.validator.Struct(i); err != nil {
    // Optionally return the error to let each route control the status code.
    return echo.ErrBadRequest.Wrap(err)
  }
  return nil
}
func main() {
  dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {

		dsn = "host=localhost user=postgres password=123456 dbname=ticketbooking port=5432 sslmode=disable TimeZone=Asia/Dhaka"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		panic("Failed to connect to database");
	} else{
    fmt.Println("Database connection successful");
	}

	
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("Failed to migrate database:" + err.Error())
	}
	fmt.Println("Database migrated successfully");


  e := echo.New()
  e.Validator = &CustomValidator{validator: validator.New()}

  e.Use(middleware.RequestLogger())
  e.Use(middleware.Recover())

  e.GET("/", func(c *echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
  });


e.POST("/users", func(c *echo.Context) error {
	newUser:=new(User);

	 if err := c.Bind(newUser); err != nil {
      return err
    }
    if err := c.Validate(newUser); err != nil {
      return err
    }
  // Handle POST request for creating a new user
  result:=db.Create(newUser);
  if result.Error != nil {
    return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create user"})
  }
  return c.JSON(http.StatusCreated, map[string]interface{}{"message": "User created successfully","data": newUser})
});




  if err := e.Start(":1323"); err != nil {
    e.Logger.Error("failed to start server", "error", err)
  }
}



