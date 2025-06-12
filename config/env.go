package config

import "os"

type EnvConfig struct {
	AuthSecret            string
	DatabaseConfig        *DatabaseConfig
	ExternalServiceConfig *ExternalServiceConfig
	AWSConfig             *AWSConfig
}

type ExternalServiceConfig struct {
	ReplicateModel    string
	ReplicateAPIKey   string
	NewsAPIURL        string
	NewsAPIKey        string
	PolyMarketBaseURL string
}

type AWSConfig struct {
	SESAccessKey       string
	SESSecretAccessKey string
	Region             string
}

func NewEnvConfig() *EnvConfig {
	return &EnvConfig{
		AuthSecret: os.Getenv("AUTH_SECRET"),
		DatabaseConfig: &DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
		},
		ExternalServiceConfig: &ExternalServiceConfig{
			ReplicateModel:    os.Getenv("REPLICATE_MODEL"),
			ReplicateAPIKey:   os.Getenv("REPLICATE_KEY"),
			NewsAPIURL:        os.Getenv("NEWS_API_URL"),
			NewsAPIKey:        os.Getenv("NEWS_API_KEY"),
			PolyMarketBaseURL: os.Getenv("POLYMARKET_BASE_URL"),
		},
		AWSConfig: &AWSConfig{
			SESAccessKey:       os.Getenv("AWS_SES_ACCESS_KEY"),
			SESSecretAccessKey: os.Getenv("AWS_SES_SECRET_ACCESS_KEY"),
			Region:             os.Getenv("AWS_REGION"),
		},
	}
}
