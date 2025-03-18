package registry

import (
	"github.com/jmoiron/sqlx"
	"github.com/qoentz/evedict/config"
	"github.com/qoentz/evedict/internal/db/repository"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"github.com/qoentz/evedict/internal/eventfeed/polymarket"
	"github.com/qoentz/evedict/internal/llm/replicate"
	"github.com/qoentz/evedict/internal/service"
)

type Registry struct {
	ForecastService   *service.ForecastService
	ReplicateService  *replicate.Service
	NewsAPIService    *newsapi.Service
	PolyMarketService *polymarket.Service
}

func NewRegistry(c *config.SystemConfig, db *sqlx.DB) *Registry {
	forecastRepository := repository.NewForecastRepository(db)

	replicateService := replicate.NewReplicateService(c.HTTPClient, c.PromptTemplate, c.EnvConfig.ExternalServiceConfig.ReplicateModel, c.EnvConfig.ExternalServiceConfig.ReplicateAPIKey)
	newsAPIService := newsapi.NewNewsAPIService(c.HTTPClient, c.EnvConfig.ExternalServiceConfig.NewsAPIKey, c.EnvConfig.ExternalServiceConfig.NewsAPIURL)

	polyMarketService := polymarket.NewPolyMarketService(c.HTTPClient, c.EnvConfig.ExternalServiceConfig.PolyMarketBaseURL)

	marketService := service.NewMarketService(polyMarketService, replicateService)
	forecastService := service.NewForecastService(forecastRepository, replicateService, newsAPIService, marketService)

	return &Registry{
		ForecastService:   forecastService,
		ReplicateService:  replicateService,
		NewsAPIService:    newsAPIService,
		PolyMarketService: polyMarketService,
	}
}
