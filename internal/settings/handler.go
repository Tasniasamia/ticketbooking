package settings

import (
	
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/settings/dto"
	"github.com/labstack/echo/v5"
)
type handler struct{ service Service }
func NewHandler(svc Service) *handler { return &handler{service: svc} }
func (h *handler) Get(c *echo.Context) error {
	res, err := h.service.Get()
	if err != nil {
		return c.JSON(500, httpresponse.Error{Success: false, StatusCode: 500, Error: true, ErrorMessage: err.Error()})
	}
	return c.JSON(200, httpresponse.Success{Success: true, StatusCode: 200, Message: "Settings fetched successfully", Data: res})
}
func (h *handler) Upsert(c *echo.Context) error {
	var req dto.UpsertRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, httpresponse.Error{Success: false, StatusCode: 400, Error: true, ErrorMessage: err.Error()})
	}
	res, err := h.service.Upsert(&req)
	if err != nil {
		return c.JSON(500, httpresponse.Error{Success: false, StatusCode: 500, Error: true, ErrorMessage: err.Error()})
	}
	return c.JSON(200, httpresponse.Success{Success: true, StatusCode: 200, Message: "Settings saved successfully", Data: res})
}
