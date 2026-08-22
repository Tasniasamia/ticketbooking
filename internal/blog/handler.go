package blog

import (
	"net/http"
	"strconv"

	"ticketBooking/internal/blog/dto"
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/utils/query"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *Service
}

func NewHandler(s *Service) *handler {
	return &handler{service: s}
}

func (h *handler) CreateBlog(c *echo.Context) error {
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

	if req.AuthorID == 0 {
		if uid, ok := c.Get("user_id").(uint); ok {
			req.AuthorID = uid
		}
	}

	res, err := h.service.CreateBlog(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to create blog", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success: true, StatusCode: http.StatusCreated,
		Message: "Blog created successfully", Data: res,
	})
}

func (h *handler) GetAllBlogs(c *echo.Context) error {
	params := query.Parse(c)
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	data, err := h.service.GetAllBlogs(params, lang)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to fetch blogs", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Blogs fetched successfully", Data: data,
	})
}

func (h *handler) GetAllBlogsPublic(c *echo.Context) error {
	params := query.Parse(c)
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}
	params.Filters["status"] = "approved"

	data, err := h.service.GetAllBlogs(params, lang)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to fetch blogs", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Blogs fetched successfully", Data: data,
	})
}

func (h *handler) GetBlogByID(c *echo.Context) error {
	blogID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Invalid blog ID",
		})
	}

	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	// optional auth — is_liked দেখানোর জন্য
	var userID uint
	if uid, ok := c.Get("user_id").(uint); ok {
		userID = uid
	}

	res, err := h.service.GetBlogByID(uint(blogID), lang, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to fetch blog", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Blog fetched successfully", Data: res,
	})
}

func (h *handler) GetBlogByIDAdmin(c *echo.Context) error {
	blogID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Invalid blog ID",
		})
	}

	res, err := h.service.GetBlogByIDAdmin(uint(blogID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to fetch blog", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Blog fetched successfully", Data: res,
	})
}

func (h *handler) UpdateBlog(c *echo.Context) error {
	blogID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Invalid blog ID",
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

	res, err := h.service.UpdateBlog(uint(blogID), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to update blog", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Blog updated successfully", Data: res,
	})
}

func (h *handler) DeleteBlog(c *echo.Context) error {
	blogID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Invalid blog ID",
		})
	}

	if err := h.service.DeleteBlog(uint(blogID)); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to delete blog", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Blog deleted successfully",
	})
}

func (h *handler) UpdateBlogStatus(c *echo.Context) error {
	var req dto.UpdateBlogStatusRequest
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

	if err := h.service.UpdateBlogStatus(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to update blog status", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Blog status updated successfully",
	})
}

func (h *handler) GetMyBlogs(c *echo.Context) error {
	authorID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success: false, StatusCode: 401, Error: true, ErrorMessage: "unauthorized",
		})
	}
	params := query.Parse(c)
	params.Filters["author_id"] = authorID

	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	data, err := h.service.GetAllBlogs(params, lang)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to fetch blogs", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Blogs fetched successfully", Data: data,
	})
}

// POST /blogs/:id/like  — toggle like (auth required)
func (h *handler) ToggleLike(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success: false, StatusCode: http.StatusUnauthorized, Error: true,
			ErrorMessage: "unauthorized",
		})
	}

	blogID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Invalid blog ID",
		})
	}

	res, err := h.service.ToggleLike(uint(blogID), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Failed to toggle like", ErrorDetails: err.Error(),
		})
	}

	msg := "Blog liked"
	if !res.Liked {
		msg = "Blog unliked"
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: msg, Data: res,
	})
}
