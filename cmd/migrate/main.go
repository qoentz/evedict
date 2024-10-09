package main

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/qoentz/evedict/config"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	envConfig := config.NewEnvConfig()

	m, err := migrate.New(
		"file://internal/db/migrations",
		envConfig.DatabaseConfig.ConfigureDSN(config.URIFormat))
	if err != nil {
		log.Fatalf("Failed to initialize migrate: %v", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}
