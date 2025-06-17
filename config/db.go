package config

import (
	"fmt"
	"log"
)

type DatabaseConfig struct {
	Host     string `env:"DB_HOST,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Name     string `env:"DB_NAME,required"`
	Port     string `env:"DB_PORT,required"`
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
