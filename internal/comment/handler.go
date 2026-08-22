package comment

import (
	"net/http"
	"strconv"

	"ticketBooking/internal/comment/dto"
	"ticketBooking/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *Service
}

func NewHandler(s *Service) *handler {
	return &handler{service: s}
}

func (h *handler) CreateComment(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success: false, StatusCode: http.StatusUnauthorized, Error: true,
			ErrorMessage: "unauthorized",
		})
	}

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

	res, err := h.service.CreateComment(userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Failed to create comment", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success: true, StatusCode: http.StatusCreated,
		Message: "Comment created successfully", Data: res,
	})
}

// GET /comments/blog/:blogId  → nested tree
func (h *handler) GetCommentsByBlog(c *echo.Context) error {
	blogID, err := strconv.ParseUint(c.Param("blogId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Invalid blog ID",
		})
	}

	tree, err := h.service.GetCommentsTree(uint(blogID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to fetch comments", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Comments fetched successfully", Data: tree,
	})
}

func (h *handler) UpdateComment(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success: false, StatusCode: http.StatusUnauthorized, Error: true,
			ErrorMessage: "unauthorized",
		})
	}

	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Invalid comment ID",
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

	role, _ := c.Get("user_role").(string)
	isAdmin := role == "admin"

	res, err := h.service.UpdateComment(uint(commentID), userID, &req, isAdmin)
	if err != nil {
		return c.JSON(http.StatusForbidden, httpresponse.Error{
			Success: false, StatusCode: http.StatusForbidden, Error: true,
			ErrorMessage: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Comment updated successfully", Data: res,
	})
}

func (h *handler) DeleteComment(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success: false, StatusCode: http.StatusUnauthorized, Error: true,
			ErrorMessage: "unauthorized",
		})
	}

	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Invalid comment ID",
		})
	}

	role, _ := c.Get("user_role").(string)
	isAdmin := role == "admin"

	if err := h.service.DeleteComment(uint(commentID), userID, isAdmin); err != nil {
		return c.JSON(http.StatusForbidden, httpresponse.Error{
			Success: false, StatusCode: http.StatusForbidden, Error: true,
			ErrorMessage: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Comment deleted successfully",
	})
}
