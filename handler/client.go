package handler

import (
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// pingPeriodRatio is the ratio for calculating ping period from pong wait
	pingPeriodRatio = 9

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * pingPeriodRatio) / 10

	// channelBufferSize is the buffer size for client send channel
	channelBufferSize = 256
)

// Client represents a WebSocket client
type Client struct {
	ID   string
	Conn *websocket.Conn
	Room *Room
	Send chan []byte
}

// NewClient creates a new client
func NewClient(id string, conn *websocket.Conn, room *Room) *Client {
	return &Client{
		ID:   id,
		Conn: conn,
		Room: room,
		Send: make(chan []byte, channelBufferSize),
	}
}

// ReadPump pumps messages from the WebSocket connection to the room
func (c *Client) ReadPump() {
	defer func() {
		c.Room.Unregister <- c
		_ = c.Conn.Close()
	}()

	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Parse command: "COMMAND ARG1 ARG2 ..."
		msg := strings.TrimSpace(string(message))
		parts := strings.Fields(msg)
		if len(parts) == 0 {
			continue
		}

		cmd := parts[0]
		args := parts[1:]

		c.Room.HandleMessage(c, cmd, args)
	}
}

// WritePump pumps messages from the room to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The room closed the channel
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Add queued messages to the current WebSocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				_, _ = w.Write([]byte{'\n'})
				_, _ = w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
