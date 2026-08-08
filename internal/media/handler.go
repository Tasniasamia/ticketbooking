package media

import (
	"net/http"
	"strconv"

	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/media/dto"
	"ticketBooking/internal/utils/query"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *Service
}

func NewHandler(s *Service) *handler {
	return &handler{service: s}
}

// UploadFile
// POST /api/v1/media/upload
// Content-Type: multipart/form-data
//
// Form fields expected from frontend:
//
//	file         (required)  – the binary file
//	model_name   (required)  – event | user | product | ...
//	model_id     (optional)  – related entity id
//	folder       (optional)  – cloudinary folder override
func (h *handler) UploadFile(c *echo.Context) error {
	// 1. Parse multipart
	file, header, err := c.Request().FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusBadRequest,
			Error:        true,
			ErrorMessage: "file is required (multipart field name must be 'file')",
			ErrorDetails: err.Error(),
		})
	}
	defer file.Close()

	// 2. Bind extra form fields
	var req dto.UploadRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusBadRequest,
			Error:        true,
			ErrorMessage: "Invalid form fields",
			ErrorDetails: err.Error(),
		})
	}

	// 3. Validate
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusBadRequest,
			Error:        true,
			ErrorMessage: "Validation failed",
			ErrorDetails: err.Error(),
		})
	}
	// ModelName string `form:"model_name" validate:"required,oneof=event user product booking category banner"`
	// ModelID   *uint  `form:"model_id"`
	// Folder    string `form:"folder"`
	// 4. Service
	res, err := h.service.UploadFile(c.Request().Context(), file, header, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusInternalServerError,
			Error:        true,
			ErrorMessage: "Failed to upload file",
			ErrorDetails: err.Error(),
		})
	}

	// 5. Success – also include the model-keyed shape the user asked for
	// keyed := map[string]interface{}{
	// 	req.ModelName + "_url": res.SecureURL,
	// 	req.ModelName + "_id":  res.ID,
	// 	"image_id":             res.ImageID,
	// 	"public_url":           res.PublicURL,
	// 	"type":                 res.Type,
	// 	"model_name":           res.ModelName,
	// 	"created_at":           res.CreatedAt,
	// }

	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success:    true,
		StatusCode: http.StatusCreated,
		Message:    "File uploaded successfully",
		Data: map[string]interface{}{
			"media": res,
			// "keyed": keyed, // e.g. event_url + event_id
		},
	})
}

// GetByID
// GET /api/v1/media/:id
func (h *handler) GetByID(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Invalid media id",
		})
	}

	res, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.Error{
			Success: false, StatusCode: http.StatusNotFound, Error: true,
			ErrorMessage: "Media not found", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Media fetched successfully", Data: res,
	})
}

// List
// GET /api/v1/media?model_name=event&model_id=5&type=image&page=1&limit=10
func (h *handler) List(c *echo.Context) error {
	params := query.Parse(c)

	filter := dto.ListFilter{
		ModelName: c.QueryParam("model_name"),
		Type:      c.QueryParam("type"),
	}
	if mid := c.QueryParam("model_id"); mid != "" {
		if v, err := strconv.ParseUint(mid, 10, 64); err == nil {
			id := uint(v)
			filter.ModelID = &id
		}
	}

	data, err := h.service.List(params, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to list media", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Media list fetched successfully", Data: data,
	})
}

// Delete
// DELETE /api/v1/media/:id
// Uses the auto-generated DB id.
// Internally finds the row → takes image_id (Cloudinary public_id) → destroys on Cloudinary → soft-deletes DB row.
func (h *handler) Delete(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest, Error: true,
			ErrorMessage: "Invalid media id",
		})
	}

	if err := h.service.DeleteByID(c.Request().Context(), uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
			ErrorMessage: "Failed to delete media", ErrorDetails: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Media deleted successfully (Cloudinary + database)",
	})
}
