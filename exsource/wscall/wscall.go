package wscall

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"

	"binance-orderbook-average/dependency"
	"binance-orderbook-average/types"
)

type WsCall struct {
	ReadMessageFromStream func() (<-chan types.BinanceDepthUpdate, error)
}

func New(d *dependency.Dependency) *WsCall {
	return &WsCall{
		ReadMessageFromStream: readMessagesFromStream(d.BinanceWS),
	}
}

var readMsgOnce sync.Once

func readMessagesFromStream(conn *websocket.Conn) func() (<-chan types.BinanceDepthUpdate, error) {
	return func() (<-chan types.BinanceDepthUpdate, error) {
		depthUpdates := make(chan types.BinanceDepthUpdate)

		readMsgOnce.Do(func() {
			go func() {
				defer close(depthUpdates)
				for {
					_, message, err := conn.ReadMessage()
					if err != nil {
						log.Println("Binance read error:", err)
						break
					}

					var depthUpdate types.BinanceDepthUpdate
					err = json.Unmarshal(message, &depthUpdate)
					if err != nil {
						log.Println("Unmarshal error:", err)
						continue
					}

					depthUpdates <- depthUpdate
				}
			}()
		})

		return depthUpdates, nil
	}
}
