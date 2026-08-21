package payment

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/payment/dto"
	"ticketBooking/internal/utils/query"

	"github.com/labstack/echo/v5"
)

type handler struct{ service Service }

func NewHandler(svc Service) *handler { return &handler{service: svc} }

func (h *handler) CreateCheckout(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{Success: false, StatusCode: 401, Error: true, ErrorMessage: "unauthorized"})
	}
	userName, _ := c.Get("user_name").(string)
	userEmail, _ := c.Get("user_email").(string)
	var req dto.CreateCheckoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{Success: false, StatusCode: 400, Error: true, ErrorMessage: "Invalid request body", ErrorDetails: err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{Success: false, StatusCode: 400, Error: true, ErrorMessage: "Validation failed", ErrorDetails: err.Error()})
	}
	res, err := h.service.CreateCheckout(userID, userName, userEmail, req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case ErrGatewayDisabled, ErrInvalidMethod, ErrMissingGatewayKey, ErrSettingsNotFound:
			status = http.StatusBadRequest
		}
		return c.JSON(status, httpresponse.Error{Success: false, StatusCode: status, Error: true, ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{Success: true, StatusCode: 200, Message: "Checkout session created", Data: res})
}

func (h *handler) GetPayment(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{Success: false, StatusCode: 400, Error: true, ErrorMessage: "Invalid payment id"})
	}
	res, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.Error{Success: false, StatusCode: 404, Error: true, ErrorMessage: "Payment not found"})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{Success: true, StatusCode: 200, Message: "Payment fetched", Data: res})
}

func (h *handler) GetByTransactionID(c *echo.Context) error {
	res, err := h.service.GetByTransactionID(c.Param("transaction_id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.Error{Success: false, StatusCode: 404, Error: true, ErrorMessage: "Payment not found"})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{Success: true, StatusCode: 200, Message: "Payment fetched", Data: res})
}

func (h *handler) GetMyPayments(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{Success: false, StatusCode: 401, Error: true, ErrorMessage: "unauthorized"})
	}
	res, err := h.service.GetByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{Success: false, StatusCode: 500, Error: true, ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{Success: true, StatusCode: 200, Message: "Payments fetched", Data: res})
}

func (h *handler) StripeWebhook(c *echo.Context) error {
	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot read body"})
	}

	sig := c.Request().Header.Get("Stripe-Signature")
	
	// Debug log
	fmt.Println("=== WEBHOOK DEBUG ===")
	fmt.Println("Signature:", sig)
	fmt.Println("Payload length:", len(payload))
	
	if err := h.service.HandleStripeWebhook(payload, sig); err != nil {
		fmt.Println("Webhook Error:", err.Error())   // ← এই লাইনটা খুব জরুরি
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"received": "true"})
}

func (h *handler) SSLCommerzIPN(c *echo.Context) error {
	var ipn dto.SSLCommerzIPN
	if err := c.Bind(&ipn); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if ipn.TranID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "tran_id required"})
	}
	if err := h.service.HandleSSLCommerzIPN(ipn); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.String(http.StatusOK, "VALID")
}
func (h *handler) VerifySSLCommerzSession(c *echo.Context) error {
	tranID := c.QueryParam("tran_id")
	status := c.QueryParam("status")

	if tranID == "" {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: 400, Error: true,
			ErrorMessage: "tran_id is required",
		})
	}

	res, err := h.service.VerifySSLCommerzSession(tranID, status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: 400, Error: true,
			ErrorMessage: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: 200,
		Message: "Payment verified",
		Data:    res,
	})
}

func (h *handler) GetAllPayments(c *echo.Context) error {
	lang := c.QueryParam("lang")
     
	if lang == "" {
		lang = "en"
	}
	params := query.Parse(c)
    res, err := h.service.GetAllPayments(params,lang)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{Success: false, StatusCode: 500, Error: true, ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{Success: true, StatusCode: 200, Message: "Payments fetched", Data: res})
}

func (h *handler) GetUserPayments(c *echo.Context) error {
    userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{Success: false, StatusCode: 401, Error: true, ErrorMessage: "unauthorized User"})
	}
			lang := c.QueryParam("lang")
     
	if lang == "" {
		lang = "en"
	}
    params := query.Parse(c)
    res, err := h.service.GetUserPayments(params, userID,lang)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{Success: false, StatusCode: 500, Error: true, ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{Success: true, StatusCode: 200, Message: "Payments fetched", Data: res})
}

func (h *handler) GetManagementPayments(c *echo.Context) error {
    managerId, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{Success: false, StatusCode: 401, Error: true, ErrorMessage: "unauthorized Manager"})
	}
		lang := c.QueryParam("lang")
     
	if lang == "" {
		lang = "en"
	}

    params := query.Parse(c)
    res, err := h.service.GetManagerPayments(params,managerId,lang)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{Success: false, StatusCode: 500, Error: true, ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{Success: true, StatusCode: 200, Message: "Payments fetched", Data: res})
}









// ---- Payment Method handlers ----

func (h *handler) CreatePaymentMethod(c *echo.Context) error {
	var req dto.CreatePaymentMethodRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{Success: false, StatusCode: 400, Error: true, ErrorMessage: "Invalid request body", ErrorDetails: err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{Success: false, StatusCode: 400, Error: true, ErrorMessage: "Validation failed", ErrorDetails: err.Error()})
	}
	res, err := h.service.CreatePaymentMethod(req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case ErrInvalidMethod, ErrPaymentMethodExists:
			status = http.StatusBadRequest
		}
		return c.JSON(status, httpresponse.Error{Success: false, StatusCode: status, Error: true, ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusCreated, httpresponse.Success{Success: true, StatusCode: 201, Message: "Payment method created", Data: res})
}

func (h *handler) UpdatePaymentMethod(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{Success: false, StatusCode: 400, Error: true, ErrorMessage: "Invalid id"})
	}
	var req dto.UpdatePaymentMethodRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{Success: false, StatusCode: 400, Error: true, ErrorMessage: "Invalid request body", ErrorDetails: err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{Success: false, StatusCode: 400, Error: true, ErrorMessage: "Validation failed", ErrorDetails: err.Error()})
	}
	res, err := h.service.UpdatePaymentMethod(uint(id), req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case ErrPaymentMethodNotFound:
			status = http.StatusNotFound
		case ErrInvalidMethod, ErrPaymentMethodExists:
			status = http.StatusBadRequest
		}
		return c.JSON(status, httpresponse.Error{Success: false, StatusCode: status, Error: true, ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{Success: true, StatusCode: 200, Message: "Payment method updated", Data: res})
}

func (h *handler) DeletePaymentMethod(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{Success: false, StatusCode: 400, Error: true, ErrorMessage: "Invalid id"})
	}
	if err := h.service.DeletePaymentMethod(uint(id)); err != nil {
		status := http.StatusInternalServerError
		if err == ErrPaymentMethodNotFound {
			status = http.StatusNotFound
		}
		return c.JSON(status, httpresponse.Error{Success: false, StatusCode: status, Error: true, ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{Success: true, StatusCode: 200, Message: "Payment method deleted", Data: nil})
}

func (h *handler) GetPaymentMethod(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{Success: false, StatusCode: 400, Error: true, ErrorMessage: "Invalid id"})
	}
	res, err := h.service.GetPaymentMethod(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.Error{Success: false, StatusCode: 404, Error: true, ErrorMessage: "Payment method not found"})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{Success: true, StatusCode: 200, Message: "Payment method fetched", Data: res})
}

func (h *handler) ListPaymentMethods(c *echo.Context) error {
	// Public list: only enabled by default. Admin can pass ?all=true
	enabledOnly := true
	if c.QueryParam("all") == "true" {
		enabledOnly = false
	}
	res, err := h.service.ListPaymentMethods(enabledOnly)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{Success: false, StatusCode: 500, Error: true, ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, httpresponse.Success{Success: true, StatusCode: 200, Message: "Payment methods fetched", Data: res})
}

func (h *handler) VerifyStripeSession(c *echo.Context) error {
	sessionID := c.QueryParam("session_id")
	if sessionID == "" {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: 400, Error: true,
			ErrorMessage: "session_id is required",
		})
	}

	res, err := h.service.VerifyStripeSession(sessionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false, StatusCode: 400, Error: true,
			ErrorMessage: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: 200,
		Message: "Payment verified",
		Data:    res,
	})
}