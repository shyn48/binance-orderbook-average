package core

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"binance-orderbook-average/clientmanager"
	"binance-orderbook-average/exsource"
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

func BroadcastAveragePrice(es *exsource.Exsource, cm *clientmanager.ClientManager) {
	depthUpdates, err := es.WsCall.ReadMessageFromStream()
	if err != nil {
		log.Fatalf("Error starting Binance stream: %v", err)
	}

	for depthUpdate := range depthUpdates {
		// Calculate the average price
		averagePrice, err := calculateAveragePrice(depthUpdate.Bids, depthUpdate.Asks)
		if err != nil {
			log.Println("Error calculating average price:", err)
			continue
		}

		message := clientmanager.Message{
			Type: "averagePrice",
			Data: averagePrice,
		}

		seralizedMessage, err := json.Marshal(message)
		if err != nil {
			log.Println("Error serializing message:", err)
			continue
		}

		// Send the average price to the broadcast channel
		cm.Broadcast <- seralizedMessage
	}
}
