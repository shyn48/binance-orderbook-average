package dependency

import (
	"fmt"

	"github.com/gorilla/websocket"

	"binance-orderbook-average/common"
)

func newBinanceConnection(cfg common.Config) (*websocket.Conn, func(), error) {
	wsURL := fmt.Sprintf("%s/%s@depth%d@%dms", cfg.Binance.Endpoint, cfg.Binance.Symbol, cfg.Binance.Depth, cfg.Binance.UpdateSpeed)

	c, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		if resp != nil {
			common.Logger().Infof("Binance connection HTTP status: %s", resp.Status)
		}
		return nil, nil, fmt.Errorf("binance WebSocket connection error: %w", err)
	}

	common.Logger().Info("Connected to Binance WebSocket")

	return c, func() { c.Close() }, nil
}
