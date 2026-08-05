package user

import (
	"net/http"

	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/user/dto"
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

func (h *handler) CreateUser(c *echo.Context) error {
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
	res, err := h.service.CreateUser(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusInternalServerError,
			Error:        true,
			ErrorMessage: "Failed to create user",
			ErrorDetails: err.Error(),
		})
	}

	// Success
	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success:    true,
		StatusCode: http.StatusCreated,
		Message:    "User created successfully",
		Data:       res,
	})
}

func (h *handler) LoginUser(c *echo.Context) error {
	var req dto.LoginRequest

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
	res, err := h.service.LoginUser(&req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusUnauthorized,
			Error:        true,
			ErrorMessage: "Invalid email or password",
			ErrorDetails: err.Error(),
		})
	}

	// Success
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "Login successful",
		Data:       res,
	})
}


func (h *handler) GetMe(c *echo.Context) error {

	userId, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusUnauthorized,
			Error:        true,
			ErrorMessage: "Unauthorized",
			ErrorDetails: "Cannot get user information from authentication middleware",
		})
	}


    
	res,err := h.service.GetMe(userId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusInternalServerError,
			Error:        true,
			ErrorMessage: "Failed to retrieve user information",
			ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "User information retrieved successfully",
		Data:       res,
	})
}