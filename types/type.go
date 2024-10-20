package types

type BinanceDepthUpdate struct {
	LastUpdateID int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
