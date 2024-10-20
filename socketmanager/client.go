package socketmanager

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"binance-orderbook-average/common"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	manager  *Manager
	conn     *websocket.Conn
	send     chan []byte
	mu       sync.Mutex
	isClosed bool
}

func (c *Client) readPump() {
	defer func() {
		c.manager.unregister <- c
		c.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				common.Logger().Errorf("client readPump error: %v", err)
			}
			break
		}
		// Ignore messages from clients (if any)
	}
}

func (c *Client) writePump() {
	// ticker for ping messages
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// batch send the rest of the buffer messages if any
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) Send(message []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isClosed {
		return fmt.Errorf("client is closed")
	}
	select {
	case c.send <- message:
		return nil
	default:
		c.Close()
		return fmt.Errorf("client send buffer full")
	}
}

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isClosed {
		c.conn.Close()
		close(c.send)
		c.isClosed = true
	}
}
