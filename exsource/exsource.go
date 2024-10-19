package exsource

import (
	"binance-orderbook-average/dependency"
	"binance-orderbook-average/exsource/wscall"
)

type Exsource struct {
	WsCall *wscall.WsCall
}

func New(d *dependency.Dependency) *Exsource {
	return &Exsource{
		WsCall: wscall.New(d),
	}
}
