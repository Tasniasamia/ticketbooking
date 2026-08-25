package messaging

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"ticketBooking/internal/auth"
	"ticketBooking/internal/httpresponse"

	"ticketBooking/internal/messaging/dto"

	"ticketBooking/internal/messaging/ws"

	"github.com/gorilla/websocket"

	"github.com/labstack/echo/v5"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // production-এ specific origin দাও
	},
}

type Handler struct {
	svc        *Service
	hub        *Hub
	jwtService auth.JwtService
}

func NewHandler(svc *Service, hub *Hub, jwtService auth.JwtService) *Handler {
	return &Handler{svc: svc, hub: hub, jwtService: jwtService}
}

// ---------- REST endpoints ----------

func (h *Handler) StartConversation(c *echo.Context) error {
	var userID uint;
	var req dto.StartConversationRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body", err)
	}
	if err := c.Validate(&req); err != nil {
		return badRequest(c, "validation failed", err)
	}
	if(req.UserId != 0) {
		 userID = req.UserId;
	}

	res, err := h.svc.StartConversation(userID, &req)
	if err != nil {
		if err == ErrCannotMessage {
			return c.JSON(http.StatusForbidden, httpresponse.Error{
				Success: false, StatusCode: 403, Error: true, ErrorMessage: err.Error(),
			})
		}
		return serverError(c, err)
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: 200, Message: "conversation ready", Data: res,
	})
}

func (h *Handler) SendMessage(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return unauthorized(c)
	}

	var req dto.SendMessageRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body", err)
	}
	if err := c.Validate(&req); err != nil {
		return badRequest(c, "validation failed", err)
	}

	res, err := h.svc.SendMessage(userID, &req)
	if err != nil {
		if err == ErrCannotMessage {
			return c.JSON(http.StatusForbidden, httpresponse.Error{
				Success: false, StatusCode: 403, Error: true, ErrorMessage: err.Error(),
			})
		}
		return serverError(c, err)
	}

	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success: true, StatusCode: 201, Message: "message sent", Data: res,
	})
}

func (h *Handler) GetMyConversations(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return unauthorized(c)
	}

	list, err := h.svc.GetMyConversations(userID)
	if err != nil {
		return serverError(c, err)
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: 200, Message: "conversations fetched", Data: list,
	})
}

func (h *Handler) GetMessages(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return unauthorized(c)
	}

	convID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return badRequest(c, "invalid conversation id", err)
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))

	msgs, total, err := h.svc.GetMessages(userID, uint(convID), page, pageSize)
	if err != nil {
		if err == ErrForbidden || err == ErrConversationNotFound {
			return c.JSON(http.StatusForbidden, httpresponse.Error{
				Success: false, StatusCode: 403, Error: true, ErrorMessage: err.Error(),
			})
		}
		return serverError(c, err)
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: 200, Message: "messages fetched",
		Data: map[string]any{
			"messages": msgs,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

func (h *Handler) MarkRead(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return unauthorized(c)
	}

	convID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return badRequest(c, "invalid conversation id", err)
	}

	var req dto.MarkReadRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body", err)
	}

	if err := h.svc.MarkRead(userID, uint(convID), req.MessageIDs); err != nil {
		if err == ErrForbidden {
			return c.JSON(http.StatusForbidden, httpresponse.Error{
				Success: false, StatusCode: 403, Error: true, ErrorMessage: err.Error(),
			})
		}
		return serverError(c, err)
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true, StatusCode: 200, Message: "messages marked as read",
	})
}

// ---------- WebSocket endpoint ----------

// GET /api/v1/messaging/ws?token=<jwt>
func (h *Handler) ServeWS(c *echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		// also accept Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}
	}
	if token == "" {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success: false, StatusCode: 401, Error: true, ErrorMessage: "missing token",
		})
	}

	claims, err := h.jwtService.ValidateToken(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success: false, StatusCode: 401, Error: true, ErrorMessage: "invalid token",
		})
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("ws upgrade error:", err)
		return nil
	}

	client := &ws.Client{
		UserID: claims.UserID,
		Send:   make(chan []byte, 256),
		Hub:    h.hub,
	}
	
	h.hub.Register(client)

	// writer
	go func() {
		defer func() {
			h.hub.Unregister(client)
			conn.Close()
		}()
		for msg := range client.Send {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		}
	}()

	// reader
	go func() {
		defer func() {
			h.hub.Unregister(client)
			conn.Close()
		}()

		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var incoming dto.WSIncoming
			if err := json.Unmarshal(data, &incoming); err != nil {
				h.sendWSError(client, "invalid json")
				continue
			}

			switch incoming.Type {
			case "ping":
				h.hub.SendToUser(client.UserID, dto.WSOutgoing{Type: "pong", Payload: nil})

			case "message":
				var payload dto.WSMessagePayload
				raw, _ := json.Marshal(incoming.Payload)
				if err := json.Unmarshal(raw, &payload); err != nil {
					h.sendWSError(client, "invalid message payload")
					continue
				}
				if payload.Content == "" || payload.ReceiverID == 0 {
					h.sendWSError(client, "receiver_id and content required")
					continue
				}
				if _, err := h.svc.HandleWSMessage(client.UserID, payload); err != nil {
					h.sendWSError(client, err.Error())
				}

			case "read":
				var payload dto.WSReadPayload
				raw, _ := json.Marshal(incoming.Payload)
				if err := json.Unmarshal(raw, &payload); err != nil {
					h.sendWSError(client, "invalid read payload")
					continue
				}
				if err := h.svc.HandleWSRead(client.UserID, payload); err != nil {
					h.sendWSError(client, err.Error())
				}

			default:
				h.sendWSError(client, "unknown type: "+incoming.Type)
			}
		}
	}()

	return nil
}

func (h *Handler) sendWSError(client *ws.Client, msg string) {
	h.hub.SendToUser(client.UserID, dto.WSOutgoing{
		Type:    "error",
		Payload: map[string]string{"message": msg},
	})
}

// ---------- helpers ----------

func unauthorized(c *echo.Context) error {
	return c.JSON(http.StatusUnauthorized, httpresponse.Error{
		Success: false, StatusCode: 401, Error: true, ErrorMessage: "unauthorized",
	})
}

func badRequest(c *echo.Context, msg string, err error) error {
	details := ""
	if err != nil {
		details = err.Error()
	}
	return c.JSON(http.StatusBadRequest, httpresponse.Error{
		Success: false, StatusCode: 400, Error: true, ErrorMessage: msg, ErrorDetails: details,
	})
}

func serverError(c *echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, httpresponse.Error{
		Success: false, StatusCode: 500, Error: true, ErrorMessage: err.Error(),
	})
}
