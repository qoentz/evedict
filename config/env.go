package config

import "os"

type EnvConfig struct {
	DatabaseConfig        *DatabaseConfig
	ExternalServiceConfig *ExternalServiceConfig
}

type ExternalServiceConfig struct {
	ReplicateModel  string
	ReplicateAPIKey string
	NewsAPIURL      string
	NewsAPIKey      string
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
