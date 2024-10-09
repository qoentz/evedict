package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

const (
	maxConnectionAttempts = 3
	initialBackoff        = 2 * time.Second
)

func InitDB(dsn string) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	backoff := initialBackoff
	for attempts := 0; attempts < maxConnectionAttempts; attempts++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err != nil {
			log.Printf("Failed to connect to database, attempt %d: %v", attempts+1, err)
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		if err = db.Ping(); err != nil {
			log.Printf("Database connection established but validation failed, attempt %d: %v", attempts+1, err)
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		log.Println("Database connection established")
		break

	}

	if err != nil {
		return nil, fmt.Errorf("failed to establish database connection: %v", err)
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	return db, nil
}
