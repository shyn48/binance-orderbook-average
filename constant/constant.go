package constant

import "time"

const (
	WebSocketReadBufferSize  = 4096
	WebSocketWriteBufferSize = 4096

	ThrottlingInterval = 300 * time.Millisecond

	WorkerCount = 500
	TaskBuffer  = 100000

	ManagerMaxActiveConn = 10000
)
