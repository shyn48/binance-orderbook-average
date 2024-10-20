package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"binance-orderbook-average/api"
	"binance-orderbook-average/common"
	"binance-orderbook-average/dependency"
	"binance-orderbook-average/exsource"
	"binance-orderbook-average/price"
	"binance-orderbook-average/socketmanager"
	"binance-orderbook-average/workerpool"
)

func main() {
	config, err := common.GetConfig()
	if err != nil {
		common.Logger().Fatal(err)
	}

	d, gracefulClose, err := dependency.New(config)
	if err != nil {
		common.Logger().Fatal(err)
	}
	defer gracefulClose()

	es := exsource.New(d)
	if err != nil {
		common.Logger().Fatal(err)
	}

	workerCount := 500
	taskBuffer := 100000
	managerWorkerPool := workerpool.New(workerCount, taskBuffer)

	common.Logger().Info("Initializing manager")
	manager := socketmanager.New(10000, managerWorkerPool)
	common.Logger().Info("Manager initialized")

	common.Logger().Info("Starting to broadcast average prices")
	go price.BroadcastAveragePrice(es, manager)
	common.Logger().Info("Started broadcasting average prices")

	api.Setup(manager)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		common.Logger().Info("Shutting down server...")

		manager.Shutdown()

		common.Logger().Info("Server gracefully stopped.")
		os.Exit(0)
	}()

	common.Logger().Printf("Starting server on port %d", config.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
