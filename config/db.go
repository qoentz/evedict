package config

import (
	"fmt"
	"log"
)

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

type DSNFormat int

const (
	URIFormat DSNFormat = iota
	KeyValueFormat
)

func (e *DatabaseConfig) ConfigureDSN(format DSNFormat) string {
	switch format {
	case URIFormat:
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			e.User, e.Password, e.Host, e.Port, e.Name)
	case KeyValueFormat:
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			e.Host, e.User, e.Password, e.Name, e.Port)
	default:
		log.Fatalf("Unknown DSN format")
		return ""
	}
}
