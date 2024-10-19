package api

import (
	"fmt"
	"net/http"

	"binance-orderbook-average/clientmanager"
)

func Setup(cm *clientmanager.ClientManager) {
	// Root route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "HelloWorld")
	})

	// WebSocket route for average price
	http.HandleFunc("/average-price", cm.HandleConnections)
}
