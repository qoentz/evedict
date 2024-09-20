package registry

import (
	"evedict/config"
	"evedict/internal/llm/replicate"
	"evedict/internal/source/newsapi"
)

type Registry struct {
	ReplicateService *replicate.Service
	NewsAPIService   *newsapi.Service
}

func NewRegistry(c *config.SystemConfig) *Registry {
	replicateService := replicate.NewReplicateService(c.HTTPClient, c.EnvConfig.ReplicateModel, c.EnvConfig.ReplicateAPIKey)
	newsAPIService := newsapi.NewNewsAPIService(c.EnvConfig.NewsAPIKey, c.EnvConfig.NewsAPIURL)
	return &Registry{
		ReplicateService: replicateService,
		NewsAPIService:   newsAPIService,
	}
}
