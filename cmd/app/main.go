package main

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/qoentz/evedict/config"
	"github.com/qoentz/evedict/internal/db"
	"github.com/qoentz/evedict/internal/registry"
	"github.com/qoentz/evedict/internal/server"
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

	database, err := db.InitDB(systemConfig.EnvConfig.DatabaseConfig.ConfigureDSN(config.KeyValueFormat))
	if err != nil {
		log.Fatalf("Error initilizing database: %v", err)
	}
	defer func(db *sqlx.DB) { _ = db.Close() }(database)

	reg := registry.NewRegistry(systemConfig, database)

	httpServer := server.ServeHTTP(server.InitRouter(reg))

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
}
