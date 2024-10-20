package price

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"binance-orderbook-average/common"
	"binance-orderbook-average/constant"
	"binance-orderbook-average/exsource"
	"binance-orderbook-average/socketmanager"
	"binance-orderbook-average/types"
)

func calculateAveragePrice(bids, asks [][]string) (float64, error) {
	var totalSum float64
	var count int

	for _, bid := range bids {
		price, err := strconv.ParseFloat(bid[0], 64)
		if err != nil {
			return 0, err
		}
		totalSum += price
		count++
	}

	for _, ask := range asks {
		price, err := strconv.ParseFloat(ask[0], 64)
		if err != nil {
			return 0, err
		}
		totalSum += price
		count++
	}

	if count == 0 {
		return 0, fmt.Errorf("no bids or asks to calculate average price")
	}

	averagePrice := totalSum / float64(count)
	return averagePrice, nil
}

func BroadcastAveragePrice(es *exsource.Exsource, sm *socketmanager.Manager) {
	depthUpdates, err := es.WsCall.ReadMessageFromStream()
	if err != nil {
		common.Logger().Fatalf("error starting Binance stream: %v", err)
	}

	ticker := time.NewTicker(constant.ThrottlingInterval)
	defer ticker.Stop()

	var lastDepthUpdate *types.BinanceDepthUpdate

	for {
		select {
		case depthUpdate, ok := <-depthUpdates:
			if !ok {
				return
			}

			lastDepthUpdate = &depthUpdate
		case <-ticker.C:
			if lastDepthUpdate != nil {
				averagePrice, err := calculateAveragePrice(lastDepthUpdate.Bids, lastDepthUpdate.Asks)
				if err != nil {
					common.Logger().Error("error calculating average price:", err)
					continue
				}

				message := types.Message{
					Type: "averagePrice",
					Data: averagePrice,
				}

				serializedMessage, err := json.Marshal(&message)
				if err != nil {
					common.Logger().Error("error serializing message:", err)
					continue
				}

				sm.BroadcastMsg(serializedMessage)

				lastDepthUpdate = nil
			}
		}
	}
}
