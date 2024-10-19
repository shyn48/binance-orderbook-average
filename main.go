// main.go

package main

import (
	"fmt"
	"log"
	"net/http"

	"binance-orderbook-average/api"
	"binance-orderbook-average/clientmanager"
	"binance-orderbook-average/common"
	"binance-orderbook-average/core"
	"binance-orderbook-average/dependency"
	"binance-orderbook-average/exsource"
)

func main() {
	// Load the configuration
	config, err := common.GetConfig()
	if err != nil {
		common.Logger().Fatal(err)
	}

	// Initialize dependencies (Binance WebSocket connection)
	d, gracefulClose, err := dependency.New(config)
	if err != nil {
		common.Logger().Fatal(err)
	}
	defer gracefulClose()

	// Initialize WebSocket call handler
	es := exsource.New(d)
	if err != nil {
		common.Logger().Fatal(err)
	}

	// Initialize the client manager
	fmt.Println("Initializing client manager")
	cm := clientmanager.NewClientManager()
	fmt.Println("Client manager initialized")

	// Setup API routes
	fmt.Println("Setting up API routes")
	api.Setup(cm)
	fmt.Println("API routes set up")

	fmt.Println("Starting client manager")
	go cm.StartBroadcasting()
	fmt.Println("Client manager started")

	// Start broadcasting average prices
	fmt.Println("Starting to broadcast average prices")
	go core.BroadcastAveragePrice(es, cm)
	fmt.Println("Started broadcasting average prices")

	// Start the HTTP server
	common.Logger().Infof("Starting server on port %d", config.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
