package config

import (
	"evedict/internal/promptgen"
	"fmt"
	"github.com/joho/godotenv"
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
		HTTPClient:     client}, nil
}

func NewEnvConfig() *EnvConfig {
	replicateModel := os.Getenv("REPLICATE_MODEL")
	replicateKey := os.Getenv("REPLICATE_KEY")
	newsAPIUrl := os.Getenv("NEWS_API_URL")
	newsAPIKey := os.Getenv("NEWS_API_KEY")

	return &EnvConfig{
		ReplicateModel:  replicateModel,
		ReplicateAPIKey: replicateKey,
		NewsAPIURL:      newsAPIUrl,
		NewsAPIKey:      newsAPIKey,
	}

}
