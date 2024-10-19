package dependency

import (
	"github.com/gorilla/websocket"

	"binance-orderbook-average/common"
)

type Dependency struct {
	BinanceWS *websocket.Conn
}

func cleanupFunctions(fns ...func()) func() {
	return func() {
		for _, fn := range fns {
			fn()
		}
	}
}

func New(cfg common.Config) (*Dependency, func(), error) {
	binanceWS, binanceCleanup, err := newBinanceConnection(cfg)
	if err != nil {
		return nil, nil, err
	}

	cleanup := cleanupFunctions(binanceCleanup)

	return &Dependency{
		BinanceWS: binanceWS,
	}, cleanup, nil
}
