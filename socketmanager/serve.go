package socketmanager

import (
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"

	"binance-orderbook-average/common"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: true,
}

func ServeWs(manager *Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt64(&manager.activeConnCount) >= manager.maxActiveConn {
			http.Error(w, "Server is busy. Please try again later.", http.StatusServiceUnavailable)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			common.Logger().Error("serveWs: Upgrade error:", err)
			return
		}

		client := &Client{
			manager: manager,
			conn:    conn,
			send:    make(chan []byte, 256),
		}

		manager.register <- client

		go client.writePump()
		go client.readPump()
	}
}
