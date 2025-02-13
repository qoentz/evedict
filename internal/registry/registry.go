package registry

import (
	"github.com/jmoiron/sqlx"
	"github.com/qoentz/evedict/config"
	"github.com/qoentz/evedict/internal/db/repository"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"github.com/qoentz/evedict/internal/llm/replicate"
	"github.com/qoentz/evedict/internal/service"
)

type Registry struct {
	DivinationService *service.DivinationService
	ReplicateService  *replicate.Service
	NewsAPIService    *newsapi.Service
}

func NewRegistry(c *config.SystemConfig, db *sqlx.DB) *Registry {
	divinationRepository := repository.NewDivinationRepository(db)

	replicateService := replicate.NewReplicateService(c.HTTPClient, c.EnvConfig.ExternalServiceConfig.ReplicateModel, c.EnvConfig.ExternalServiceConfig.ReplicateAPIKey)
	newsAPIService := newsapi.NewNewsAPIService(c.HTTPClient, c.EnvConfig.ExternalServiceConfig.NewsAPIKey, c.EnvConfig.ExternalServiceConfig.NewsAPIURL)

	divinationService := service.NewDivinationService(divinationRepository, replicateService, newsAPIService, c.PromptTemplate)

	return &Registry{
		DivinationService: divinationService,
		ReplicateService:  replicateService,
		NewsAPIService:    newsAPIService,
	}
}
