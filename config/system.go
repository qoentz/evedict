package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/qoentz/evedict/internal/promptgen"
	"log"
	"net/http"
	"os"
	"time"
)

type SystemConfig struct {
	EnvConfig      *EnvConfig
	PromptTemplate *promptgen.PromptTemplate
	HTTPClient     *http.Client
}

type EnvConfig struct {
	DatabaseConfig        *DatabaseConfig
	ExternalServiceConfig *ExternalServiceConfig
}

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

type ExternalServiceConfig struct {
	ReplicateModel  string
	ReplicateAPIKey string
	NewsAPIURL      string
	NewsAPIKey      string
}

func ConfigureSystem() (*SystemConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	envConfig := NewEnvConfig()

	promptTemplate, err := promptgen.LoadPromptTemplate("internal/promptgen/prompts.yaml")
	if err != nil {
		return nil, fmt.Errorf("error loading prompt template: %v", err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &SystemConfig{
		EnvConfig:      envConfig,
		PromptTemplate: promptTemplate,
		HTTPClient:     client,
	}, nil
}

func NewEnvConfig() *EnvConfig {
	return &EnvConfig{
		DatabaseConfig: &DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
		},
		ExternalServiceConfig: &ExternalServiceConfig{
			ReplicateModel:  os.Getenv("REPLICATE_MODEL"),
			ReplicateAPIKey: os.Getenv("REPLICATE_KEY"),
			NewsAPIURL:      os.Getenv("NEWS_API_URL"),
			NewsAPIKey:      os.Getenv("NEWS_API_KEY"),
		},
	}
}

func (e *DatabaseConfig) ConfigureDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		e.Host, e.User, e.Password, e.Name, e.Port)
}
