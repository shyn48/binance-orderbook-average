package socketmanager

import (
	"sync/atomic"

	"github.com/gorilla/websocket"

	"binance-orderbook-average/common"
	"binance-orderbook-average/workerpool"
)

type ManagerTask struct {
	Client  *Client
	Message []byte
}

type Manager struct {
	BroadCastMsg func([]byte)
	Shutdown     func()

	processTask func(workerpool.Task)

	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client

	workerPool *workerpool.WorkerPool

	activeConnCount int64
	maxActiveConn   int64

	quit chan struct{}
}

func New(maxActiveConn int64, workerPool *workerpool.WorkerPool) *Manager {
	manager := &Manager{
		clients:       make(map[*Client]bool),
		broadcast:     make(chan []byte),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		workerPool:    workerPool,
		maxActiveConn: maxActiveConn,
		quit:          make(chan struct{}),
	}

	manager.BroadCastMsg = broadcastMsg(manager)
	manager.processTask = processTask(manager)
	manager.Shutdown = shutdown(manager)

	go run(manager)

	manager.workerPool.Start(manager.processTask)

	return manager
}

func broadcastMsg(m *Manager) func(message []byte) {
	return func(message []byte) {
		m.broadcast <- message
	}
}

func run(m *Manager) {
	for {
		select {
		case client := <-m.register:
			if atomic.LoadInt64(&m.activeConnCount) >= m.maxActiveConn {
				// Exceeded max connections, close the client
				client.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseTryAgainLater, "Server is busy"))
				client.Close()
				continue
			}
			m.clients[client] = true
			atomic.AddInt64(&m.activeConnCount, 1)
			common.Logger().Tracef("Run: Client registered. Active connections: %d", m.activeConnCount)
		case client := <-m.unregister:
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				client.Close()
				atomic.AddInt64(&m.activeConnCount, -1)
			}
			common.Logger().Tracef("Run: Client unregistered. Active connections: %d", m.activeConnCount)
		case message := <-m.broadcast:
			for client := range m.clients {
				task := ManagerTask{
					Client:  client,
					Message: message,
				}

				m.workerPool.Submit(task)
			}
		case <-m.quit:
			common.Logger().Infof("Run: Received quit signal. Shutting down Manager.")
			for client := range m.clients {
				client.Close()
				delete(m.clients, client)
				atomic.AddInt64(&m.activeConnCount, -1)
			}
			return

		}
	}
}

func processTask(m *Manager) func(task workerpool.Task) {
	return func(task workerpool.Task) {
		sendTask, ok := task.(ManagerTask)
		if !ok {
			common.Logger().Error("processTask: invalid task type")
			return
		}

		err := sendTask.Client.Send(sendTask.Message)
		if err != nil {
			m.unregister <- sendTask.Client
			atomic.AddInt64(&m.activeConnCount, -1)
			common.Logger().Warnf("processTask: Removed client %s due to %v", sendTask.Client.conn.RemoteAddr(), err)
		} else {
			common.Logger().Tracef("processTask: Message sent to client %s", sendTask.Client.conn.RemoteAddr())
		}
	}
}

func shutdown(m *Manager) func() {
	return func() {
		close(m.quit)
		m.workerPool.Stop()
	}
}
