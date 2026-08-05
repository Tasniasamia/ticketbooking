package language

import (
	"net/http"
	"strconv"
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/domain/language/dto"
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
			ErrorMessage: "Failed to create language", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success: true, StatusCode: http.StatusCreated,
		Message: "Language created successfully", Data: res,
	})
}

func (h *handler) GetAll(c *echo.Context) error {
	params := query.Parse(c)
	lang := c.QueryParam("lang")


	if lang == "" {
		lang = "en"
	}

	res, err := h.service.GetAll(params, lang)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to fetch languages", ErrorDetails: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Languages fetched successfully", Data: res,
	})
}

func (h *handler) GetByID(c *echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	res, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.Error{
			Success: false, StatusCode: http.StatusNotFound, Error: true,
			ErrorMessage: "Language not found", ErrorDetails: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Language fetched successfully", Data: res,
	})
}

func (h *handler) Update(c *echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dto.UpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Invalid request body", ErrorDetails: err.Error(),
		})
	}

	res, err := h.service.Update(uint(id), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to update language", ErrorDetails: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Language updated successfully", Data: res,
	})
}

func (h *handler) Delete(c *echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Failed to delete language", ErrorDetails: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Language deleted successfully",
	})
}