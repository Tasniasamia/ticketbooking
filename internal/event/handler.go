// package event

// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"ticketBooking/internal/event/dto"
// 	"ticketBooking/internal/httpresponse"
// 	"ticketBooking/internal/utils/query"

// 	"github.com/labstack/echo/v5"
// )

// type handler struct {
// 	service *Service
// }

// func NewHandler(s *Service) *handler {
// 	return &handler{
// 		service: s,
// 	}
// }

// func (h *handler) CreateEvent(c *echo.Context) error {
// 	var req dto.CreateRequest

// 	// Bind
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, httpresponse.Error{
// 			Success:      false,
// 			StatusCode:   http.StatusBadRequest,
// 			Error:        true,
// 			ErrorMessage: "Invalid request body",
// 			ErrorDetails: err.Error(),
// 		})
// 	}

// 	// Validate
// 	if err := c.Validate(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, httpresponse.Error{
// 			Success:      false,
// 			StatusCode:   http.StatusBadRequest,
// 			Error:        true,
// 			ErrorMessage: "Validation failed",
// 			ErrorDetails: err.Error(),
// 		})
// 	}

// 	// Service call
// 	res, err := h.service.CreateEvent(&req)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
// 			Success:      false,
// 			StatusCode:   http.StatusInternalServerError,
// 			Error:        true,
// 			ErrorMessage: "Failed to create event",
// 			ErrorDetails: err.Error(),
// 		})
// 	}

// 	// Success
// 	return c.JSON(http.StatusCreated, httpresponse.Success{
// 		Success:    true,
// 		StatusCode: http.StatusCreated,
// 		Message:    "Event created successfully",
// 		Data:       res,
// 	})
// }

// func (h *handler) GetAllEvents(c *echo.Context) error {
// 	params := query.Parse(c)
// 	lang := c.QueryParam("lang")
     
// 	if lang == "" {
// 		lang = "en"
// 	}


// 	data, err := h.service.GetAllEvents(params, lang)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
// 			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
// 			ErrorMessage: "Failed to fetch events", ErrorDetails: err.Error(),
// 		})
// 	}

// 	return c.JSON(http.StatusOK, httpresponse.Success{
// 		Success: true, StatusCode: http.StatusOK,
// 		Message: "Events fetched successfully", Data: data,
// 	})
// }
// func (h *handler) GetAllEventsPublic(c *echo.Context) error {
// 	params := query.Parse(c)
// 	lang := c.QueryParam("lang")
     
// 	if lang == "" {
// 		lang = "en"
// 	}

// 	params.Filters["status"] = "approved";

// 	data, err := h.service.GetAllEvents(params, lang)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
// 			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
// 			ErrorMessage: "Failed to fetch events", ErrorDetails: err.Error(),
// 		})
// 	}

// 	return c.JSON(http.StatusOK, httpresponse.Success{
// 		Success: true, StatusCode: http.StatusOK,
// 		Message: "Events fetched successfully", Data: data,
// 	})
// }


// func (h *handler) DeleteEvent(c *echo.Context) error {
// 	eventId, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, httpresponse.Error{
// 			Success: false, StatusCode: http.StatusBadRequest, Error: true,
// 			ErrorMessage: "Invalid event ID",
// 		})
// 	}

// 	err = h.service.DeleteEvent(uint(eventId))
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
// 			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
// 			ErrorMessage: "Failed to delete event", ErrorDetails: err.Error(),
// 		})
// 	}

// 	return c.JSON(http.StatusOK, httpresponse.Success{
// 		Success: true, StatusCode: http.StatusOK,
// 		Message: "Event deleted successfully",
// 	})
// }

// func (h *handler) UpdateEvent(c *echo.Context) error {
// 	eventId, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, httpresponse.Error{
// 			Success: false, StatusCode: http.StatusBadRequest, Error: true,
// 			ErrorMessage: "Invalid event ID",
// 		})
// 	}

// 	var req dto.UpdateRequest
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, httpresponse.Error{
// 			Success:      false,
// 			StatusCode:   http.StatusBadRequest,
// 			Error:        true,
// 			ErrorMessage: "Invalid request body",
// 			ErrorDetails: err.Error(),
// 		})
// 	}

// 	if err := c.Validate(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, httpresponse.Error{
// 			Success:      false,
// 			StatusCode:   http.StatusBadRequest,
// 			Error:        true,
// 			ErrorMessage: "Validation failed",
// 			ErrorDetails: err.Error(),
// 		})
// 	}

// 	res, err := h.service.UpdateEvent(uint(eventId), &req)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
// 			Success:      false,
// 			StatusCode:   http.StatusInternalServerError,
// 			Error:        true,
// 			ErrorMessage: "Failed to update event",
// 			ErrorDetails: err.Error(),
// 		})
// 	}

// 	return c.JSON(http.StatusOK, httpresponse.Success{
// 		Success:    true,
// 		StatusCode: http.StatusOK,
// 		Message:    "Event updated successfully",
// 		Data:       res,
// 	})
// }

// func (h *handler) GetEventByID(c *echo.Context) error {
// 	eventId, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, httpresponse.Error{
// 			Success: false, StatusCode: http.StatusBadRequest, Error: true,
// 			ErrorMessage: "Invalid event ID",
// 		})
// 	}

// 	lang := c.QueryParam("lang")
// 	if lang == "" {
// 		lang = "en"
// 	}

// 	res, err := h.service.GetEventByID(uint(eventId), lang)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
// 			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
// 			ErrorMessage: "Failed to fetch event", ErrorDetails: err.Error(),
// 		})
// 	}

// 	return c.JSON(http.StatusOK, httpresponse.Success{
// 		Success: true, StatusCode: http.StatusOK,
// 		Message: "Event fetched successfully", Data: res,
// 	})
// }

// func (h *handler) GetEventByIDAdmin(c *echo.Context) error {
// 	eventId, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, httpresponse.Error{
// 			Success: false, StatusCode: http.StatusBadRequest, Error: true,
// 			ErrorMessage: "Invalid event ID",
// 		})
// 	}

// 	res, err := h.service.GetEventByIDAdmin(uint(eventId))
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
// 			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
// 			ErrorMessage: "Failed to fetch event", ErrorDetails: err.Error(),
// 		})
// 	}

// 	return c.JSON(http.StatusOK, httpresponse.Success{
// 		Success: true, StatusCode: http.StatusOK,
// 		Message: "Event fetched successfully", Data: res,
// 	})
// }


// func (h *handler) GetMyEvents(c *echo.Context) error {
// 	managerID, ok := c.Get("user_id").(uint)
// 	if !ok {
// 		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
// 			Success: false, StatusCode: 401, Error: true, ErrorMessage: "unauthorized",
// 		})
// 	}
// 	params := query.Parse(c)
// 	params.Filters["manager_id"] = managerID

// 	lang := c.QueryParam("lang")

// 	fmt.Println("lang", lang);
// 	if lang == "" {
// 		lang = "en"
// 	}

// 	data, err := h.service.GetAllEvents(params, lang)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
// 			Success: false, StatusCode: http.StatusInternalServerError, Error: true,
// 			ErrorMessage: "Failed to fetch events", ErrorDetails: err.Error(),
// 		})
// 	}

// 	return c.JSON(http.StatusOK, httpresponse.Success{
// 		Success: true, StatusCode: http.StatusOK,
// 		Message: "Events fetched successfully", Data: data,
// 	})
// }

// func (h *handler) UpdateEventStatus(c *echo.Context) error {
// 	var req dto.UpdateEventStatusRequest
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, httpresponse.Error{
// 			Success:      false,
// 			StatusCode:   http.StatusBadRequest,
// 			Error:        true,
// 			ErrorMessage: "Invalid request body",
// 			ErrorDetails: err.Error(),
// 		})
// 	}

// 	if err := c.Validate(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, httpresponse.Error{
// 			Success:      false,
// 			StatusCode:   http.StatusBadRequest,
// 			Error:        true,
// 			ErrorMessage: "Validation failed",
// 			ErrorDetails: err.Error(),
// 		})
// 	}

// 	if err := h.service.UpdateEventStatus(&req); err != nil {
// 		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
// 			Success:      false,
// 			StatusCode:   http.StatusInternalServerError,
// 			Error:        true,
// 			ErrorMessage: "Failed to update event status",
// 			ErrorDetails: err.Error(),
// 		})
// 	}

// 	return c.JSON(http.StatusOK, httpresponse.Success{
// 		Success:    true,
// 		StatusCode: http.StatusOK,
// 		Message:    "Event status updated successfully",
// 	})
// }

package event

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"ticketBooking/internal/event/dto"
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/utils/query"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type handler struct {
	service *Service
}

func NewHandler(s *Service) *handler {
	return &handler{service: s}
}

/* ---------------- helpers ---------------- */

func fail(c *echo.Context, status int, msg string, err error) error {
	res := httpresponse.Error{
		Success:      false,
		StatusCode:   status,
		Error:        true,
		ErrorMessage: msg,
	}
	if err != nil {
		res.ErrorDetails = err.Error()
	}
	return c.JSON(status, res)
}

// serviceFail: service এর error কে ঠিক status এ বদলায়।
func serviceFail(c *echo.Context, msg string, err error) error {
	switch {
	case errors.Is(err, ErrForbidden):
		return fail(c, http.StatusForbidden, "Forbidden", err)
	case errors.Is(err, gorm.ErrRecordNotFound):
		return fail(c, http.StatusNotFound, "Event not found", err)
	default:
		return fail(c, http.StatusInternalServerError, msg, err)
	}
}

func idParam(c *echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	return uint(id), err
}

// actor: token থেকে আসা user id আর role।
func actor(c *echo.Context) (uint, string, bool) {
	id, ok := c.Get("user_id").(uint)
	// role RoleType হোক বা string, fmt.Sprint দুটোকেই string করে
	return id, fmt.Sprint(c.Get("user_role")), ok
}

func ok(c *echo.Context, status int, msg string, data any) error {
	return c.JSON(status, httpresponse.Success{
		Success:    true,
		StatusCode: status,
		Message:    msg,
		Data:       data,
	})
}

// list: তিন রকম list এর common অংশ। extra দিয়ে filter বসানো হয়।
// lang শুধু search এর জন্য (title ->> 'en'), response এ সবসময় পুরো multi-lang object যায়।
func (h *handler) list(c *echo.Context, extra func(p *query.Params)) error {
	params := query.Parse(c)
	if extra != nil {
		extra(&params)
	}

	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	data, err := h.service.GetAllEvents(params, lang)
	if err != nil {
		return fail(c, http.StatusInternalServerError, "Failed to fetch events", err)
	}
	return ok(c, http.StatusOK, "Events fetched successfully", data)
}

/* ---------------- Manager ---------------- */

// POST /events: manager তৈরি করে, event pending থাকে।
func (h *handler) CreateEvent(c *echo.Context) error {
	managerID, isOK := c.Get("user_id").(uint)
	if !isOK {
		return fail(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	var req dto.CreateRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "Invalid request body", err)
	}
	if err := c.Validate(&req); err != nil {
		return fail(c, http.StatusBadRequest, "Validation failed", err)
	}

	res, err := h.service.CreateEvent(&req, managerID)
	if err != nil {
		return fail(c, http.StatusInternalServerError, "Failed to create event", err)
	}
	return ok(c, http.StatusCreated, "Event created successfully", res)
}

// GET /events/myEvents: manager এর নিজের সব event (যেকোনো status)।
func (h *handler) GetMyEvents(c *echo.Context) error {
	managerID, isOK := c.Get("user_id").(uint)
	if !isOK {
		return fail(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	return h.list(c, func(p *query.Params) {
		p.Filters["manager_id"] = managerID // query তে manager_id দিলেও এটাই জেতে
	})
}

/* ---------------- Public ---------------- */

// GET /events আর /events/public: শুধু approved।
func (h *handler) GetAllEventsPublic(c *echo.Context) error {
	return h.list(c, func(p *query.Params) {
		p.Filters["status"] = "approved"
	})
}

// GET /events/:id: শুধু approved।
func (h *handler) GetEventByID(c *echo.Context) error {
	id, err := idParam(c)
	if err != nil {
		return fail(c, http.StatusBadRequest, "Invalid event ID", err)
	}

	res, err := h.service.GetEventByID(id, c.QueryParam("lang"))
	if err != nil {
		return serviceFail(c, "Failed to fetch event", err)
	}
	return ok(c, http.StatusOK, "Event fetched successfully", res)
}

/* ---------------- Admin ---------------- */

// GET /events/admin?status=pending&search=music&page=1: সব event।
func (h *handler) GetAllEventsAdmin(c *echo.Context) error {
	return h.list(c, nil)
}

// GET /events/admin/:id
func (h *handler) GetEventByIDAdmin(c *echo.Context) error {
	id, err := idParam(c)
	if err != nil {
		return fail(c, http.StatusBadRequest, "Invalid event ID", err)
	}

	res, err := h.service.GetEventByIDAdmin(id)
	if err != nil {
		return serviceFail(c, "Failed to fetch event", err)
	}
	return ok(c, http.StatusOK, "Event fetched successfully", res)
}

// PATCH /events/updateStatus  { "eventId": 5, "status": "approved" }
func (h *handler) UpdateEventStatus(c *echo.Context) error {
	var req dto.UpdateEventStatusRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "Invalid request body", err)
	}
	if err := c.Validate(&req); err != nil {
		return fail(c, http.StatusBadRequest, "Validation failed", err)
	}

	if err := h.service.UpdateEventStatus(&req); err != nil {
		return serviceFail(c, "Failed to update event status", err)
	}
	return ok(c, http.StatusOK, "Event status updated successfully", nil)
}

/* ---------------- Manager (নিজের) + Admin ---------------- */

// PUT /events/:id
func (h *handler) UpdateEvent(c *echo.Context) error {
	id, err := idParam(c)
	if err != nil {
		return fail(c, http.StatusBadRequest, "Invalid event ID", err)
	}
	actorID, role, isOK := actor(c)
	if !isOK {
		return fail(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	var req dto.UpdateRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "Invalid request body", err)
	}
	if err := c.Validate(&req); err != nil {
		return fail(c, http.StatusBadRequest, "Validation failed", err)
	}

	res, err := h.service.UpdateEvent(id, &req, actorID, role)
	if err != nil {
		return serviceFail(c, "Failed to update event", err)
	}
	return ok(c, http.StatusOK, "Event updated successfully", res)
}

// DELETE /events/:id
func (h *handler) DeleteEvent(c *echo.Context) error {
	id, err := idParam(c)
	if err != nil {
		return fail(c, http.StatusBadRequest, "Invalid event ID", err)
	}
	actorID, role, isOK := actor(c)
	if !isOK {
		return fail(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	if err := h.service.DeleteEvent(id, actorID, role); err != nil {
		return serviceFail(c, "Failed to delete event", err)
	}
	return ok(c, http.StatusOK, "Event deleted successfully", nil)
}
