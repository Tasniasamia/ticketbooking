package currency

import (
	"net/http"
	"strconv"
    "ticketBooking/internal/currency/dto"
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/utils/query"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// GET /api/v1/currencies
// ?enabled=true → only enabled
func (h *Handler) List(c *echo.Context) error {
		params := query.Parse(c)
	var (
		list *httpresponse.PaginatedData
		err  error
	)
	if c.QueryParam("enabled") == "true" {
		list, err = h.svc.ListEnabled(params)
	} else {
		list, err = h.svc.ListAll(params)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false, StatusCode: http.StatusInternalServerError,
			Error: true, ErrorMessage: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "OK", Data: list,
	})
}

// GET /api/v1/currencies/default
func (h *Handler) GetDefault(c *echo.Context) error {
	cur, err := h.svc.GetDefault()
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.Error{
			Success: false, StatusCode: http.StatusNotFound,
			Error: true, ErrorMessage: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "OK", Data: cur,
	})
}

// GET /api/v1/currencies/:code
func (h *Handler) GetByCode(c *echo.Context) error {
	code := c.Param("code")
	cur, err := h.svc.GetByCode(code)
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.Error{
			Success: false, StatusCode: http.StatusNotFound,
			Error: true, ErrorMessage: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "OK", Data: cur,
	})
}

// POST /api/v1/currencies
func (h *Handler) Create(c *echo.Context) error {
	var req dto.CreateCurrencyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest,
			Error: true, ErrorMessage: "invalid body", ErrorDetails: err.Error(),
		})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest,
			Error: true, ErrorMessage: "validation failed", ErrorDetails: err.Error(),
		})
	}

	cur, err := h.svc.Create(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest,
			Error: true, ErrorMessage: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success: true, StatusCode: http.StatusCreated,
		Message: "Currency created", Data: cur,
	})
}

// PUT /api/v1/currencies/:id
func (h *Handler) Update(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest,
			Error: true, ErrorMessage: "invalid id",
		})
	}

	var req dto.UpdateCurrencyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest,
			Error: true, ErrorMessage: "invalid body",
		})
	}

	cur, err := h.svc.Update(uint(id), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest,
			Error: true, ErrorMessage: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Currency updated", Data: cur,
	})
}

// POST /api/v1/currencies/set-default
// Body: { "code": "USD" }  — single click friendly
func (h *Handler) SetDefault(c *echo.Context) error {
	var req dto.SetDefaultRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest,
			Error: true, ErrorMessage: "invalid body",
		})
	}
	if req.Code == "" {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest,
			Error: true, ErrorMessage: "code is required",
		})
	}

	cur, err := h.svc.SetDefault(req.Code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest,
			Error: true, ErrorMessage: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Default currency updated", Data: cur,
	})
}

// POST /api/v1/currencies/convert
// Body: { "amount": 100, "from_code": "USD", "to_code": "BDT" }
func (h *Handler) Convert(c *echo.Context) error {
	var req dto.ConvertRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest,
			Error: true, ErrorMessage: "invalid body",
		})
	}
	if req.Amount <= 0 || req.FromCode == "" || req.ToCode == "" {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest,
			Error: true, ErrorMessage: "amount, from_code and to_code are required",
		})
	}

	result, err := h.svc.Convert(req.Amount, req.FromCode, req.ToCode)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest,
			Error: true, ErrorMessage: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "OK",
		Data: dto.ConvertResponse{
			Amount:   req.Amount,
			FromCode: req.FromCode,
			ToCode:   req.ToCode,
			Result:   result,
		},
	})
}

// DELETE /api/v1/currencies/:id
func (h *Handler) Delete(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest,
			Error: true, ErrorMessage: "invalid id",
		})
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: http.StatusBadRequest,
			Error: true, ErrorMessage: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: http.StatusOK,
		Message: "Currency deleted",
	})
}
