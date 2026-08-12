package settings

import (
	"net/http"
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/settings/dto"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service Service
}

func NewHandler(svc Service) *handler {
	return &handler{service: svc}
}

// Get returns the current site settings (public).
func (h *handler) Get(c *echo.Context) error {
	res, err := h.service.Get()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusInternalServerError,
			Error:        true,
			ErrorMessage: "Failed to fetch settings",
			ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "Settings fetched successfully",
		Data:       res,
	})
}

// Upsert creates settings if none exist, otherwise updates the existing one (auth required).
func (h *handler) Upsert(c *echo.Context) error {
	var req dto.UpsertRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusBadRequest,
			Error:        true,
			ErrorMessage: "Invalid request body",
			ErrorDetails: err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusBadRequest,
			Error:        true,
			ErrorMessage: "Validation failed",
			ErrorDetails: err.Error(),
		})
	}

	res, err := h.service.Upsert(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusInternalServerError,
			Error:        true,
			ErrorMessage: "Failed to save settings",
			ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "Settings saved successfully",
		Data:       res,
	})
}