package dashboard

import (
	"net/http"

	"ticketBooking/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc} }

func ok(c *echo.Context, msg string, data any) error {
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK, Message: msg, Data: data,
	})
}

func fail(c *echo.Context, status int, msg string, err error) error {
	res := httpresponse.Error{Success: false, StatusCode: status, Error: true, ErrorMessage: msg}
	if err != nil {
		res.ErrorDetails = err.Error()
	}
	return c.JSON(status, res)
}

// token থেকে আসা user id (AuthMiddleware এ uint হিসেবে set করা)
func currentUserID(c *echo.Context) (uint, bool) {
	id, isOK := c.Get("user_id").(uint)
	return id, isOK && id != 0
}

// GET /dashboard/admin
func (h *Handler) Admin(c *echo.Context) error {
	d, err := h.svc.Admin()
	if err != nil {
		return fail(c, http.StatusInternalServerError, "Failed to load admin dashboard", err)
	}
	return ok(c, "Admin dashboard fetched successfully", d)
}

// GET /dashboard/manager
func (h *Handler) Manager(c *echo.Context) error {
	id, found := currentUserID(c)
	if !found {
		return fail(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	d, err := h.svc.Manager(id)
	if err != nil {
		return fail(c, http.StatusInternalServerError, "Failed to load manager dashboard", err)
	}
	return ok(c, "Manager dashboard fetched successfully", d)
}

// GET /dashboard/me
func (h *Handler) User(c *echo.Context) error {
	id, found := currentUserID(c)
	if !found {
		return fail(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	d, err := h.svc.User(id)
	if err != nil {
		return fail(c, http.StatusInternalServerError, "Failed to load dashboard", err)
	}
	return ok(c, "Dashboard fetched successfully", d)
}
