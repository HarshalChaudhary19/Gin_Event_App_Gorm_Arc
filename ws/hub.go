// ws/hub.go
package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Hub struct {
	// channels
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte

	// connected clients
	Clients map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.Clients[c] = true
		case c := <-h.Unregister:
			if _, ok := h.Clients[c]; ok {
				delete(h.Clients, c)
				close(c.Send)
			}
		case msg := <-h.Broadcast:
			for c := range h.Clients {
				select {
				case c.Send <- msg:
				default:
					// drop a stalled client
					delete(h.Clients, c)
					close(c.Send)
				}
			}
		}
	}
}

// helper to broadcast structured JSON
func (h *Hub) BroadcastJSON(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Println("ws: marshal error:", err)
		return
	}
	h.Broadcast <- b
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMsgSize = 512
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// you should check origin in production
		return true
	},
}

// Client represents a connection
type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte

	// optional: attach user info (from JWT)
	UserID int
}

// ServeWS upgrades the Gin request to websocket and registers client
func (h *Hub) ServeWS(c *gin.Context) {
	// Optionally: validate token in query or cookie here before upgrade
	// token := c.Query("token") -> validate, set userId
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("ws: upgrade:", err)
		return
	}

	client := &Client{
		Hub:  h,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	h.Register <- client

	// readPump blocks (reads messages from client)
	go client.writePump()
	client.readPump() // when readPump returns it will unregister and close
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMsgSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		// we don't expect large client messages; you can handle client commands here
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			// read error (client disconnected)
			break
		}
		log.Printf("Message from Client : %s", msg)

		// ignore body — or handle client "subscribe" messages, etc.
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(msg)

			// flush queued messages
			n := len(c.Send)
			for i := 0; i < n; i++ {
				_, _ = w.Write([]byte{'\n'})
				_, _ = w.Write(<-c.Send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
