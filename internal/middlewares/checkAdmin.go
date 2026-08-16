package middleware

import (
	"fmt"
	"net/http"
	"ticketBooking/internal/user/dto"
	"ticketBooking/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

func AdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			role := c.Get("user_role")
			fmt.Println("role is ",role)
			if role != dto.ADMIN {
			return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusUnauthorized,
			Error:        true,
			ErrorMessage: "Admin access required",
			ErrorDetails: "Admin access required",
		})
			}
			return next(c)
		}
	}
}