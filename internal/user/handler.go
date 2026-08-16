package user

import (
	"net/http"
	"strconv"

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

	res, err := h.service.CreateUser(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success: true, StatusCode: http.StatusCreated,
		Message: "User created successfully. Please verify OTP sent to your email.",
		Data:    res,
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

func (h *handler) UpdateUser(c *echo.Context) error {

	userId, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusUnauthorized,
			Error:        true,
			ErrorMessage: "Unauthorized",
			ErrorDetails: "Cannot get user ID from authentication middleware",
		})
	}

	var req dto.UpdateRequest

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

	res, err := h.service.UpdateUser(userId, &req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusInternalServerError,
			Error:        true,
			ErrorMessage: "Failed to update user",
			ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "User updated successfully",
		Data:       res,
	})
}
func (h *handler) UpdateMember(c *echo.Context) error {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusBadRequest,
			Error:        true,
			ErrorMessage: "Invalid user ID",
			ErrorDetails: err.Error(),
		})
	}

	var req dto.UpdateMemberRequest

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

	res, err := h.service.UpdateMember(uint(id), &req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusInternalServerError,
			Error:        true,
			ErrorMessage: "Failed to update member",
			ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "Member updated successfully",
		Data:       res,
	})
}
func (h *handler) DeleteUser(c *echo.Context) error {
	userId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusBadRequest,
			Error:        true,
			ErrorMessage: "Invalid user ID",
			ErrorDetails: err.Error(),
		})
	}

	if err := h.service.DeleteUser(uint(userId)); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusInternalServerError,
			Error:        true,
			ErrorMessage: "Failed to delete user",
			ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "User deleted successfully",
	})


}

func (h *handler) VerifyOTP(c *echo.Context) error {
	var req dto.VerifyOTPRequest
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

	if err := h.service.VerifyOTP(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "OTP verified successfully",
	})
}

func (h *handler) ForgotPassword(c *echo.Context) error {
	var req dto.ForgotPasswordRequest
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

	if err := h.service.ForgotPassword(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "If the email exists, an OTP has been sent",
	})
}

func (h *handler) ResetPassword(c *echo.Context) error {
	var req dto.ResetPasswordRequest
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

	if err := h.service.ResetPassword(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Password reset successfully",
	})
}

func (h *handler) ResendOTP(c *echo.Context) error {
	var req dto.ResendOTPRequest
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

	if err := h.service.ResendOTP(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "OTP has been resent successfully",
	})
}