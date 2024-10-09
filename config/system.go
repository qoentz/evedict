package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/qoentz/evedict/internal/promptgen"
	"log"
	"net/http"
	"time"
)

type SystemConfig struct {
	EnvConfig      *EnvConfig
	PromptTemplate *promptgen.PromptTemplate
	HTTPClient     *http.Client
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
