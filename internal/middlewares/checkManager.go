package middleware

import (
	"fmt"
	"net/http"
	"ticketBooking/internal/user/dto"
	"ticketBooking/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

func ManagerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			role := c.Get("user_role")
			fmt.Println("role is ",role)
			if role != dto.MANAGER {
			return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusUnauthorized,
			Error:        true,
			ErrorMessage: "Manager access required",
			ErrorDetails: "Manager access required",
		})
			}
			return next(c)
		}
	}
}