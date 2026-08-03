package middleware

import (
	"fmt"
	
	"net/http"
	"strings"
	"ticketBooking/internal/auth"

	"github.com/labstack/echo/v5"
)

func AuthMiddleware(jwtService auth.JwtService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {

			authHeader :=c.Request().Header.Get("Authorization")
		if authHeader == ""{
			return c.JSON(http.StatusUnauthorized,map[string]string{
				"error":"Missing authorization header",
			})
		}
		ports :=strings.Split(authHeader," ")
		if len(ports) != 2 || ports[0] != "Bearer"{
			return c.JSON(http.StatusUnauthorized,map[string]string{
				"error":"Invalid authorization header",
			})
		}
		tokenstring :=ports[1];


		//validate token
		claims,err :=jwtService.ValidateToken(tokenstring)
		if err != nil {
			return c.JSON(http.StatusUnauthorized,map[string]string{
				"error":"Invalid or expired token",
			})
		}


		fmt.Println("claims",claims);
		//store user
		c.Set("user_id",claims.UserID);
		c.Set ("user_email",claims.Email);
		c.Set("user_name",claims.Name);


		// Implementation for authentication middleware
		return next(c)
	}
	}
}
