package booking

import (
	// "fmt"
	"fmt"
	"net/http"
	"strconv"
	"ticketBooking/internal/booking/dto"
	userDto "ticketBooking/internal/user/dto"
	"ticketBooking/internal/utils/query"

	"ticketBooking/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

// type Handler struct{ svc Service }

// func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		svc: s,
	}
}


func (h *Handler) CreateBooking(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success: false, StatusCode: 401, Error: true, ErrorMessage: "unauthorized",
		})
	}
	var req dto.CreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: 400, Error: true, ErrorMessage: "invalid request body", ErrorDetails: err.Error(),
		})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: 400, Error: true, ErrorMessage: "validation failed", ErrorDetails: err.Error(),
		})
	}
	res, err := h.svc.CreateBooking(userID, req)
	if err != nil {
		if err == ErrNotEnoughTickets {
			return c.JSON(http.StatusBadRequest, httpresponse.Error{
				Success: false, StatusCode: 400, Error: true, ErrorMessage: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: 500, Error: true, ErrorMessage: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success: true, StatusCode: 201, Message: "Ticket booked successfully", Data: res,
	})
}

func (h *Handler) GetMyBookings(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	fmt.Println("userId is ",userID);
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success: false, StatusCode: 401, Error: true, ErrorMessage: "unauthorized for invalid type",
		})
	}
	userRole:=c.Get("user_role")

	// fmt.Println("userRole", userRole)

	if userRole == "" {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success: false, StatusCode: 401, Error: true, ErrorMessage: "unauthorized",
		})
	}
	var res []*dto.Response
	var err error
	if userRole == userDto.MANAGER {
		res, err = h.svc.GetByManagerID(userID)
	} else {
		res, err = h.svc.GetByUserID(userID)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: 500, Error: true, ErrorMessage: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: 200, Message: "Bookings fetched", Data: res,
	})
}

func (h *Handler) GetByID(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: 400, Error: true, ErrorMessage: "invalid booking id",
		})
	}
	res, err := h.svc.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.Error{
			Success: false, StatusCode: 404, Error: true, ErrorMessage: "booking not found",
		})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: 200, Message: "Booking fetched", Data: res,
	})
}

func (h *Handler) GetAllBookings(c *echo.Context) error {
	p := query.Parse(c)

	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, httpresponse.Error{
	// 		Success: false, StatusCode: 400, Error: true, ErrorMessage: "invalid query parameters",
	// 	})
	// }
	
	res, err := h.svc.GetAll(p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: 500, Error: true, ErrorMessage: err.Error(),
		})
	}
	// return c.JSON(http.StatusOK, res)
		return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Bookings fetched successfully", Data: res,
	})
}