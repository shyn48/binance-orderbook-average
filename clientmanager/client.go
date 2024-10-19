package clientmanager

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type ClientManager struct {
	HandleConnections func(w http.ResponseWriter, r *http.Request)
	StartBroadcasting func()
	Broadcast         chan interface{}

	// Internal fields
	clients  map[*websocket.Conn]bool
	mu       sync.Mutex
	upgrader websocket.Upgrader
}

func handleConnections(cm *ClientManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := cm.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("HandleConnections: Upgrade error:", err)
			return
		}
		defer ws.Close()

		cm.mu.Lock()
		cm.clients[ws] = true
		cm.mu.Unlock()
		log.Println("HandleConnections: New client connected")

		for {
			_, _, err := ws.NextReader()
			if err != nil {
				cm.mu.Lock()
				delete(cm.clients, ws)
				cm.mu.Unlock()
				log.Println("HandleConnections: Client disconnected")
				break
			}
		}
	}
}

func startBroadcasting(cm *ClientManager) func() {
	return func() {
		for message := range cm.Broadcast {
			// Copy the clients slice while holding the mutex
			cm.mu.Lock()
			clients := make([]*websocket.Conn, 0, len(cm.clients))
			for client := range cm.clients {
				clients = append(clients, client)
			}
			cm.mu.Unlock()

			// Broadcast the message to all clients
			for _, client := range clients {
				err := client.WriteMessage(websocket.TextMessage, message.([]byte))
				if err != nil {
					log.Printf("StartBroadcasting: Error sending to client: %v", err)
					client.Close()
					cm.mu.Lock()
					delete(cm.clients, client)
					cm.mu.Unlock()
				}
			}
		}
	}
}

func NewClientManager() *ClientManager {
	cm := &ClientManager{
		clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan interface{}),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Adjust according to your security requirements
				return true
			},
			Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
				log.Printf("Upgrader: Error upgrading connection: %v", reason)
				log.Printf("Upgrader: Status: %d", status)
			},
		},
	}

	cm.HandleConnections = handleConnections(cm)
	cm.StartBroadcasting = startBroadcasting(cm)

	return cm
}
