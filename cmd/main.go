package main

import (
	"context"
	"evedict/config"
	"evedict/internal/registry"
	"evedict/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	systemConfig, err := config.ConfigureSystem()
	if err != nil {
		log.Fatalf("Error configuring system: %v", err)
	}

	reg := registry.NewRegistry(systemConfig)

	httpServer := server.ServeHTTP(server.InitRouter(systemConfig, reg))

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
}
