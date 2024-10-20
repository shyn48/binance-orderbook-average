package wscall

import (
	"encoding/json"

	"github.com/gorilla/websocket"

	"binance-orderbook-average/common"
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

func readMessagesFromStream(conn *websocket.Conn) func() (<-chan types.BinanceDepthUpdate, error) {
	return func() (<-chan types.BinanceDepthUpdate, error) {
		depthUpdates := make(chan types.BinanceDepthUpdate)

		go func() {
			defer close(depthUpdates)
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					common.Logger().Error("Binance read error:", err)
					break
				}

				var depthUpdate types.BinanceDepthUpdate
				err = json.Unmarshal(message, &depthUpdate)
				if err != nil {
					common.Logger().Error("Unmarshal error:", err)
					continue
				}

				depthUpdates <- depthUpdate
			}
		}()

		return depthUpdates, nil
	}
}
