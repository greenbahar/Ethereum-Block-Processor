package main

import (
	"context"
	"ethereum-parser/internal/routing"
	"ethereum-parser/internal/routing/handlers"
	ethereumRPC "ethereum-parser/internal/services/ethereum-rpc"
	"ethereum-parser/internal/services/syncwithmainnet"
	inmemory "ethereum-parser/internal/storage/in-memory"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	// simple load config
	godotenv.Load(".env")

	// setup dependencies
	storage := inmemory.NewStorage()
	ParserService := ethereumRPC.NewParser(storage)
	handler := handlers.NewHandler(ParserService)
	router := routing.SetUpRouters(handler)
	synchronizer := syncwithmainnet.NewSynchronizer(storage)

	// synchronizing with the mainnet to process the blocks
	ticker := time.NewTicker(10 * time.Second) // imagine 10 seconds is the average block generation rate on ETH
	go func() {
		synchronizer.SyncWithMainNetViaRPC(ctx)
		for range ticker.C {
			synchronizer.SyncWithMainNetViaRPC(ctx)
		}
	}()

	// start serving
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	fmt.Printf("server is running on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
