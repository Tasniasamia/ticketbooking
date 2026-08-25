package messaging

import "ticketBooking/internal/messaging/ws"

// Re-export Hub so other packages can create & start it without importing ws sub-package deeply.
type Hub = ws.Hub

func NewHub() *Hub {
	return ws.NewHub()
}
