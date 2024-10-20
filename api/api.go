package api

import (
	"fmt"
	"net/http"

	"binance-orderbook-average/socketmanager"
)

func Setup(sm *socketmanager.Manager) {
	// Root route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "HelloWorld")
	})

	// WebSocket route for average price
	http.HandleFunc("/average-price", socketmanager.ServeWs(sm))
}
