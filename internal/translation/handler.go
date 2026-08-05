package translation

import (
	"net/http"

	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/translation/dto"
	"ticketBooking/internal/utils/query"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *Service
}

func NewHandler(s *Service) *handler {
	return &handler{service: s}
}

func (h *handler) Create(c *echo.Context) error {
	var req dto.CreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Invalid request body", ErrorDetails: err.Error(),
		})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Validation failed", ErrorDetails: err.Error(),
		})
	}

	res, err := h.service.Create(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to create translation", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success: true, StatusCode: http.StatusCreated,
		Message: "Translation created successfully", Data: res,
	})
}

func (h *handler) UpdateByKey(c *echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Translation key is required",
		})
	}

	var req dto.UpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Invalid request body", ErrorDetails: err.Error(),
		})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Validation failed", ErrorDetails: err.Error(),
		})
	}

	res, err := h.service.UpdateByKey(key, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "translation key not found" {
			status = http.StatusNotFound
		}
		return c.JSON(status, httpresponse.Error{
			Success: false, StatusCode: status, Error: true,
			ErrorMessage: "Failed to update translation", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Translation updated successfully", Data: res,
	})
}

func (h *handler) BulkUpdate(c *echo.Context) error {
	var req dto.BulkUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Invalid request body", ErrorDetails: err.Error(),
		})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Validation failed", ErrorDetails: err.Error(),
		})
	}

	res, err := h.service.BulkUpdate(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to bulk update translations", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Translations updated successfully", Data: res,
	})
}

func (h *handler) GetAll(c *echo.Context) error {
	params := query.Parse(c)

	res, err := h.service.GetAll(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to fetch translations", ErrorDetails: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Translations fetched successfully", Data: res,
	})
}

func (h *handler) GetGroup(c *echo.Context) error {
	res, err := h.service.GetGroup()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to fetch translation group", ErrorDetails: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Translation group fetched successfully", Data: res,
	})
}

func (h *handler) GetByKey(c *echo.Context) error {
	key := c.Param("key")
	res, err := h.service.GetByKey(key)
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.Error{
			Success: false, StatusCode: http.StatusNotFound, Error: true,
			ErrorMessage: "Translation not found", ErrorDetails: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Translation fetched successfully", Data: res,
	})
}

func (h *handler) DeleteByKey(c *echo.Context) error {
	key := c.Param("key")
	if err := h.service.DeleteByKey(key); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "translation key not found" {
			status = http.StatusNotFound
		}
		return c.JSON(status, httpresponse.Error{
			Success: false, StatusCode: status, Error: true,
			ErrorMessage: "Failed to delete translation", ErrorDetails: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Translation deleted successfully",
	})
}