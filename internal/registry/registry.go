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
	replicateService := replicate.NewReplicateService(c.HTTPClient, c.EnvConfig.ReplicateModel, c.EnvConfig.ReplicateAPIKey)
	newsAPIService := newsapi.NewNewsAPIService(c.HTTPClient, c.EnvConfig.NewsAPIKey, c.EnvConfig.NewsAPIURL)
	return &Registry{
		ReplicateService: replicateService,
		NewsAPIService:   newsAPIService,
	}
}
