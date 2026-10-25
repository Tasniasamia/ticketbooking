package middleware

import (
	"fmt"
	"net/http"

	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/user/dto"

	"github.com/labstack/echo/v5"
)

// RolesMiddleware: দেওয়া role গুলোর যেকোনো একটা থাকলে ঢুকতে দেয়।
// AuthMiddleware এর পরে বসাতে হবে (সেটাই "user_role" set করে)।
func RolesMiddleware(allowed ...dto.RoleType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// RoleType হোক বা string, fmt.Sprint দুটোকেই string করে দেয়
			role := fmt.Sprint(c.Get("user_role"))

			for _, a := range allowed {
				if role == string(a) {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, httpresponse.Error{
				Success:      false,
				StatusCode:   http.StatusForbidden,
				Error:        true,
				ErrorMessage: "Forbidden User",
				ErrorDetails: "You do not have permission to perform this action",
			})
		}
	}
}

// AdminOrManagerMiddleware: admin অথবা manager।
func AdminOrManagerMiddleware() echo.MiddlewareFunc {
	return RolesMiddleware(dto.ADMIN, dto.MANAGER)
}