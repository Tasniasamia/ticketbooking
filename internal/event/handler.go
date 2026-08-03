package event

import (
	"net/http"

	"ticketBooking/internal/event/dto"
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/utils/query"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *Service
}

func NewHandler(s *Service) *handler {
	return &handler{
		service: s,
	}
}

func (h *handler) CreateEvent(c *echo.Context) error {
	var req dto.CreateRequest

	// Bind
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusBadRequest,
			Error:        true,
			ErrorMessage: "Invalid request body",
			ErrorDetails: err.Error(),
		})
	}

	// Validate
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusBadRequest,
			Error:        true,
			ErrorMessage: "Validation failed",
			ErrorDetails: err.Error(),
		})
	}

	// Service call
	res, err := h.service.CreateEvent(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusInternalServerError,
			Error:        true,
			ErrorMessage: "Failed to create event",
			ErrorDetails: err.Error(),
		})
	}

	// Success
	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success:    true,
		StatusCode: http.StatusCreated,
		Message:    "Event created successfully",
		Data:       res,
	})
}

func (h *handler) GetAllEvents(c *echo.Context) error {
	params := query.Parse(c)
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	data, err := h.service.GetAllEvents(params, lang)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to fetch events", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Events fetched successfully", Data: data,
	})
}