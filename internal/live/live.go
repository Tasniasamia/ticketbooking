
package live

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

const (
	pongWait   = 60 * time.Second // ei somoyer moddhe pong na ele connection dead dhora hobe
	pingPeriod = 30 * time.Second // ei por por server ping pathay
	writeWait  = 10 * time.Second
	sendBuffer = 32
)

// Message: browser ke ja pathano hoy.
type Message struct {
	Type    string `json:"type"` // "tickets_changed"
	EventID uint   `json:"event_id"`
}

type client struct {
	send chan []byte
	done chan struct{} // connection bondho hole close hoy (send channel ke kokhono close kora hoy na)
}

type Broadcaster struct {
	mu      sync.RWMutex
	clients map[*client]struct{}
}

// Default: pura app e ekta-i broadcaster.
var Default = New()

func New() *Broadcaster {
	return &Broadcaster{clients: make(map[*client]struct{})}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Ekhane shudhu public "kono event bodleche" bartha jay, kono private data na.
	// Production e specific frontend origin dite paren.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Publish: sob connected client ke janiye dey je ei event er ticket sonkha bodleche.
// Kokhono block kore na: kono client dhimi hole tar jonno bartha bad jay
// (dhimi client er queue age thekei bhora, tar refresh pending ache).
func (b *Broadcaster) Publish(eventID uint) {
	data, err := json.Marshal(Message{Type: "tickets_changed", EventID: eventID})
	if err != nil {
		return
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	for c := range b.clients {
		select {
		case c.send <- data:
		default:
		}
	}
}

// Serve: GET /api/v1/live/events. Login lagbe na, shudhu public data.
func (b *Broadcaster) Serve(c *echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("[live] upgrade error:", err)
		return nil
	}

	cl := &client{send: make(chan []byte, sendBuffer), done: make(chan struct{})}
	b.mu.Lock()
	b.clients[cl] = struct{}{}
	b.mu.Unlock()

	go b.writePump(conn, cl)
	go b.readPump(conn, cl)
	return nil
}

func (b *Broadcaster) remove(cl *client) {
	b.mu.Lock()
	delete(b.clients, cl)
	b.mu.Unlock()
}

// readPump: client kichu pathay na, ekhane shudhu connection bondho holo kina (ar pong) dhora hoy.
func (b *Broadcaster) readPump(conn *websocket.Conn, cl *client) {
	defer func() {
		b.remove(cl)
		close(cl.done)
		conn.Close()
	}()

	conn.SetReadLimit(512)
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}

// writePump: bartha pathay, ar periodically ping pathay jate dead connection dhora jay.
func (b *Broadcaster) writePump(conn *websocket.Conn, cl *client) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case msg := <-cl.send:
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				conn.Close() // readPump error pabe ar cleanup korbe
				return
			}
		case <-ticker.C:
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				conn.Close()
				return
			}
		case <-cl.done:
			return
		}
	}
}
