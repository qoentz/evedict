package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/qoentz/evedict/internal/promptgen"
	"net/http"
	"os"
	"time"
)

type SystemConfig struct {
	EnvConfig      *EnvConfig
	PromptTemplate *promptgen.PromptTemplate
	HTTPClient     *http.Client
}

func ConfigureSystem() (*SystemConfig, error) {
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load()
	}

	envConfig, err := NewEnvConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading .env: %v", err)
	}

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
