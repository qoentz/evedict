package registry

import (
	"github.com/qoentz/evedict/config"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"github.com/qoentz/evedict/internal/llm/replicate"
)

type Registry struct {
	ReplicateService *replicate.Service
	NewsAPIService   *newsapi.Service
}

func NewRegistry(c *config.SystemConfig) *Registry {
	replicateService := replicate.NewReplicateService(c.HTTPClient, c.EnvConfig.ExternalServiceConfig.ReplicateModel, c.EnvConfig.ExternalServiceConfig.ReplicateAPIKey)
	newsAPIService := newsapi.NewNewsAPIService(c.HTTPClient, c.EnvConfig.ExternalServiceConfig.NewsAPIKey, c.EnvConfig.ExternalServiceConfig.NewsAPIURL)
	return &Registry{
		ReplicateService: replicateService,
		NewsAPIService:   newsAPIService,
	}
}
