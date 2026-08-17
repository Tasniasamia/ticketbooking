package eventCategory

import (
	"fmt"
	"net/http"
	"strconv"

	"ticketBooking/internal/eventCategory/dto"
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

func (h *handler) CreateEventCategory(c *echo.Context) error {
	var req dto.CreateEventCategoryRequest

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
	res, err := h.service.CreateEventCategory(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusInternalServerError,
			Error:        true,
			ErrorMessage: "Failed to create event category",
			ErrorDetails: err.Error(),
		})
	}

	// Success
	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success:    true,
		StatusCode: http.StatusCreated,
		Message:    "Event category created successfully",
		Data:       res,
	})
}

func (h *handler) GetAllEventCategories(c *echo.Context) error {
	params := query.Parse(c)
	lang := c.QueryParam("lang")

	fmt.Println("lang", lang);
	if lang == "" {
		lang = "en"
	}

	data, err := h.service.GetAllEventCategories(params, lang)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to fetch event categories", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Event categories fetched successfully", Data: data,
	})
}


func (h *handler) GetEventCategoryByID(c *echo.Context) error {
	eventId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusBadRequest,
			Error:        true,
			ErrorMessage: "Invalid event category ID",
			ErrorDetails: err.Error(),
		})
	}

	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	res, err := h.service.GetEventCategoryByID(uint(eventId), lang)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusInternalServerError,
			Error:        true,
			ErrorMessage: "Failed to fetch event category",
			ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Event category fetched successfully", Data: res,
	})
}

func (h *handler) UpdateEventCategory(c *echo.Context) error {
	eventId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusBadRequest,
			Error:        true,
			ErrorMessage: "Invalid event category ID",
			ErrorDetails: err.Error(),
		})
	}

	var req dto.UpdateEventCategoryRequest
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

	res, err := h.service.UpdateEventCategory(uint(eventId), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusInternalServerError,
			Error:        true,
			ErrorMessage: "Failed to update event category",
			ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Event category updated successfully", Data: res,
	})
}


func (h *handler) DeleteEventCategory(c *echo.Context) error {
	eventCategoryId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusBadRequest,
			Error:        true,
			ErrorMessage: "Invalid event category ID",
			ErrorDetails: err.Error(),
		})
	}

	if err := h.service.DeleteEventCategory(uint(eventCategoryId)); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusInternalServerError,
			Error:        true,
			ErrorMessage: "Failed to delete event category",
			ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "Event category deleted successfully",
	})


}