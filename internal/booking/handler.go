package booking

import (
	"net/http"

	"ticketBooking/internal/booking/dto"
	"ticketBooking/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	svc *service
}

func NewHandler(svc *service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateBooking(c *echo.Context) error {

	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusUnauthorized,
			Error:        true,
			ErrorMessage: "unauthorized",
			ErrorDetails: "User not found",
		})
	}

	var req dto.CreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusBadRequest,
			Error:        true,
			ErrorMessage: "invalid request body",
			ErrorDetails: err.Error(),
		})
	}

	if req.EventID == 0 || req.Quantity < 1 {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success:      false,
			StatusCode:   http.StatusBadRequest,
			Error:        true,
			ErrorMessage: "event_id and quantity are required",
			ErrorDetails: "event_id must be greater than 0 and quantity must be at least 1",
		})
	}

	res, err := h.svc.CreateBooking(userID, req)
	if err != nil {
		switch err {
		case ErrNotEnoughTickets:
			return c.JSON(http.StatusBadRequest, httpresponse.Error{
				Success:      false,
				StatusCode:   http.StatusBadRequest,
				Error:        true,
				ErrorMessage: err.Error(),
				ErrorDetails: "Not enough tickets available",
			})
		default:
			return c.JSON(http.StatusInternalServerError, httpresponse.Error{
				Success:      false,
				StatusCode:   http.StatusInternalServerError,
				Error:        true,
				ErrorMessage: err.Error(),
				ErrorDetails: "Failed to create booking",
			})
		}
	}

	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success:    true,
		StatusCode: http.StatusCreated,
		Message:    "Ticket Booking Sucessfully",
		Data:       res,
	})

}
