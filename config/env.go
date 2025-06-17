package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type EnvConfig struct {
	AuthSecret            string `env:"AUTH_SECRET,required"`
	DatabaseConfig        *DatabaseConfig
	ExternalServiceConfig *ExternalServiceConfig
	AWSConfig             *AWSConfig
}

type ExternalServiceConfig struct {
	ReplicateModel    string `env:"REPLICATE_MODEL,required"`
	ReplicateAPIKey   string `env:"REPLICATE_KEY,required"`
	NewsAPIURL        string `env:"NEWS_API_URL,required"`
	NewsAPIKey        string `env:"NEWS_API_KEY,required"`
	PolyMarketBaseURL string `env:"POLYMARKET_BASE_URL,required"`
}

type AWSConfig struct {
	SESAccessKey       string `env:"AWS_SES_ACCESS_KEY,required"`
	SESSecretAccessKey string `env:"AWS_SES_SECRET_ACCESS_KEY,required"`
	Region             string `env:"AWS_REGION,required"`
}

func NewEnvConfig() (*EnvConfig, error) {
	config := &EnvConfig{
		DatabaseConfig:        &DatabaseConfig{},
		ExternalServiceConfig: &ExternalServiceConfig{},
		AWSConfig:             &AWSConfig{},
	}

	if err := loadEnvVars(config); err != nil {
		return nil, err
	}

	return config, nil
}

func loadEnvVars(config interface{}) error {
	v := reflect.ValueOf(config).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			// Handle nested structs recursively
			if field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct {
				if err := loadEnvVars(field.Interface()); err != nil {
					return err
				}
			}
			continue
		}

		parts := strings.Split(envTag, ",")
		envKey := parts[0]
		required := len(parts) > 1 && parts[1] == "required"

		value := os.Getenv(envKey)
		if required && value == "" {
			return fmt.Errorf("required environment variable %s is not set", envKey)
		}

		field.SetString(value)
	}

	return nil
}
