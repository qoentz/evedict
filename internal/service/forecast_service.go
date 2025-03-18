package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/qoentz/evedict/internal/api/dto"
	"github.com/qoentz/evedict/internal/db/model"
	"github.com/qoentz/evedict/internal/db/repository"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"github.com/qoentz/evedict/internal/eventfeed/polymarket"
	"github.com/qoentz/evedict/internal/llm"
	"github.com/qoentz/evedict/internal/llm/replicate"
	"github.com/qoentz/evedict/internal/promptgen"
	"github.com/qoentz/evedict/internal/util"
	"log"
	"time"
)

type ForecastService struct {
	ForecastRepository *repository.ForecastRepository
	AIService          llm.Service
	NewsAPIService     *newsapi.Service
	MarketService      *MarketService
}

func NewForecastService(forecastRepository *repository.ForecastRepository, replicateService *replicate.Service, newsAPIService *newsapi.Service, marketService *MarketService) *ForecastService {
	return &ForecastService{
		ForecastRepository: forecastRepository,
		AIService:          replicateService,
		NewsAPIService:     newsAPIService,
		MarketService:      marketService,
	}
}

func (s *ForecastService) GeneratePolyForecasts() ([]dto.Forecast, error) {
	selectedEvents, err := s.MarketService.GetMarketEvents(2)
	if err != nil {
		return nil, err
	}

	var forecasts []dto.Forecast
	for _, e := range selectedEvents {
		var keywords []string
		for _, tag := range e.Tags {
			keywords = append(keywords, tag.Label)
		}

		articles, err := s.NewsAPIService.FetchWithKeywords(keywords)
		if err != nil {
			return nil, fmt.Errorf("error fetching articles from NewsAPI: %v", err)
		}

		mainArticleIdx, err := s.AIService.SelectIndex(promptgen.SelectArticleForEvent, struct {
			Event    polymarket.Event
			Articles []newsapi.Article
		}{Event: e, Articles: articles})
		if err != nil {
			log.Printf("Error selecting article for event %s: %v", e.Title, err)
			continue
		}

		if mainArticleIdx < 0 || mainArticleIdx >= len(articles) {
			log.Printf("Invalid article index (%d) for event: %s", mainArticleIdx, e.Title)
			continue
		}

		mainArticle := articles[mainArticleIdx]

		if exists, _ := s.ForecastRepository.CheckImageURL(mainArticle.URLToImage); exists {
			log.Println(mainArticle.Title + " already exists!")
			continue
		}

		forecast, err := s.AIService.GetForecast(mainArticle, articles, &e)
		if err != nil {
			log.Printf("Error generating forecast for event %s: %v", e.Title, err)
			continue
		}

		s.MarketService.AttachMarketData(e, forecast)
		s.attachMetadata(mainArticle, forecast, keywords, articles)

		forecasts = append(forecasts, *forecast)
	}

	return forecasts, nil
}

func (s *ForecastService) GenerateForecasts(category newsapi.Category) ([]dto.Forecast, error) {
	headlines, err := s.NewsAPIService.FetchTopHeadlines(category)
	if err != nil {
		return nil, fmt.Errorf("error fetching headlines from NewsAPI: %v", err)
	}

	articleSelection, err := s.AIService.SelectIndexes(promptgen.SelectArticles, headlines, 2)
	if err != nil {
		return nil, fmt.Errorf("error selecting markets: %v", err)
	}

	var forecasts []dto.Forecast
	for _, idx := range articleSelection {
		mainArticle := headlines[idx]

		exists, _ := s.ForecastRepository.CheckImageURL(mainArticle.URLToImage)
		if exists {
			log.Println(mainArticle.Title + " already exists!")
			continue
		}

		keywords, err := s.AIService.ExtractKeywords(mainArticle)
		if err != nil {
			return nil, fmt.Errorf("error extracting keywords: %v", err)
		}

		articles, err := s.NewsAPIService.FetchWithKeywords(keywords)
		if err != nil {
			return nil, fmt.Errorf("error fetching articles from NewsAPI with keywords: %v", err)
		}

		forecast, err := s.AIService.GetForecast(mainArticle, articles, nil)
		if err != nil {
			return nil, fmt.Errorf("error generating forecast: %v", err)
		}

		s.attachMetadata(mainArticle, forecast, keywords, articles)
		forecasts = append(forecasts, *forecast)
	}

	return forecasts, nil
}

func (s *ForecastService) attachMetadata(mainArticle newsapi.Article, forecast *dto.Forecast, keywords []string, articles []newsapi.Article) {
	forecast.ImageURL = mainArticle.URLToImage
	forecast.Timestamp = time.Now().UTC()

	forecast.Tags = make([]dto.Tag, len(keywords))
	for i, keyword := range keywords {
		forecast.Tags[i] = dto.Tag{Name: keyword}
	}

	var sources []dto.Source
	sources = append(sources, dto.Source{
		Name:     mainArticle.Source.Name,
		Title:    mainArticle.Title,
		URL:      mainArticle.URL,
		ImageURL: mainArticle.URLToImage,
	})

	for _, article := range articles {
		if article.URL == mainArticle.URL || article.Title == "[Removed]" {
			continue
		}
		sources = append(sources, dto.Source{
			Name:     article.Source.Name,
			Title:    article.Title,
			URL:      article.URL,
			ImageURL: article.URLToImage,
		})
	}

	forecast.Sources = sources
}

func (s *ForecastService) GetForecasts(limit int, offset int, category *util.Category) ([]dto.Forecast, error) {
	forecasts, err := s.ForecastRepository.GetForecasts(limit, offset, category)
	if err != nil {
		return nil, fmt.Errorf("failed to get forecasts: %v", err)
	}

	var result []dto.Forecast
	for _, forecast := range forecasts {
		dtoForecast := s.convertToDTO(&forecast)
		result = append(result, *dtoForecast)
	}

	return result, nil
}

func (s *ForecastService) GetForecast(forecastID uuid.UUID) (*dto.Forecast, error) {
	forecast, err := s.ForecastRepository.GetForecast(forecastID)
	if err != nil {
		return nil, err
	}

	if forecast == nil {
		return nil, fmt.Errorf("forecast not found")
	}

	tagNames := make([]string, len(forecast.Tags))
	for i, t := range forecast.Tags {
		tagNames[i] = t.Name
	}

	relatedForecasts, err := s.ForecastRepository.GetRelatedForecastsByTagAndCategory(forecast.ID, tagNames, forecast.Category, 4)
	if err != nil {
		return nil, err
	}

	dtoForecast := s.convertToDTO(forecast)

	dtoRelated := make([]dto.Forecast, len(relatedForecasts))
	for i, rf := range relatedForecasts {
		dtoRelated[i] = dto.Forecast{
			ID:       rf.ID,
			Headline: rf.Headline,
			Summary:  rf.Summary,
			ImageURL: rf.ImageURL,
		}
	}

	dtoForecast.Related = dtoRelated

	return dtoForecast, nil
}

func (s *ForecastService) SavePolyForecasts(forecasts []dto.Forecast) error {
	modelForecasts := s.convertToModel(forecasts)
	err := s.ForecastRepository.SavePolyForecasts(modelForecasts)
	if err != nil {
		return fmt.Errorf("failed to save forecasts: %v", err)
	}
	return nil
}

func (s *ForecastService) SaveForecasts(forecasts []dto.Forecast) error {
	modelForecasts := s.convertToModel(forecasts)
	err := s.ForecastRepository.SaveForecasts(modelForecasts)
	if err != nil {
		return fmt.Errorf("failed to save forecasts: %v", err)
	}
	return nil
}

func (s *ForecastService) SaveForecast(forecast *dto.Forecast) error {
	modelForecasts := s.convertToModel([]dto.Forecast{*forecast})
	err := s.ForecastRepository.SaveForecast(&modelForecasts[0])
	if err != nil {
		return fmt.Errorf("failed to save forecast: %v", err)
	}
	return nil
}

func (s *ForecastService) convertToDTO(forecast *model.Forecast) *dto.Forecast {
	dtoOutcomes := make([]dto.Outcome, len(forecast.Outcomes))
	for i, o := range forecast.Outcomes {
		dtoOutcomes[i] = dto.Outcome{
			Content:         o.Content,
			ConfidenceLevel: o.ConfidenceLevel,
		}
	}

	dtoTags := make([]dto.Tag, len(forecast.Tags))
	for i, t := range forecast.Tags {
		dtoTags[i] = dto.Tag{
			Name: t.Name,
		}
	}

	dtoSources := make([]dto.Source, len(forecast.Sources))
	for i, src := range forecast.Sources {
		dtoSources[i] = dto.Source{
			Name:  src.Name,
			Title: src.Title,
			URL:   src.URL,
		}

		if src.ImageURL != nil {
			dtoSources[i].ImageURL = *src.ImageURL
		}
	}

	var dtoMarket *dto.Market
	if forecast.Market != nil {
		dtoMarket = &dto.Market{
			ID:            forecast.Market.ID,
			Question:      forecast.Market.Question,
			Outcomes:      forecast.Market.Outcomes,
			OutcomePrices: forecast.Market.OutcomePrices,
			Volume:        forecast.Market.Volume,
			ImageURL:      forecast.Market.ImageURL,
		}

		_ = ParseOutcomesAndPrices(dtoMarket)
	}

	return &dto.Forecast{
		ID:        forecast.ID,
		Headline:  forecast.Headline,
		Summary:   forecast.Summary,
		Outcomes:  dtoOutcomes,
		ImageURL:  forecast.ImageURL,
		Tags:      dtoTags,
		Sources:   dtoSources,
		Timestamp: forecast.Timestamp,
		Market:    dtoMarket,
	}
}

func (s *ForecastService) convertToModel(forecasts []dto.Forecast) []model.Forecast {
	modelForecasts := make([]model.Forecast, len(forecasts))

	for i, forecast := range forecasts {
		// Generate a UUID for the forecast
		forecastID := uuid.New()

		// Generate UUIDs and set ForecastID for associated Outcomes
		outcomes := make([]model.Outcome, len(forecast.Outcomes))
		for j, o := range forecast.Outcomes {
			outcomes[j] = model.Outcome{
				ID:              uuid.New(), // New UUID for each outcome
				ForecastID:      forecastID,
				Content:         o.Content,
				ConfidenceLevel: o.ConfidenceLevel,
			}
		}

		tags := make([]model.Tag, len(forecast.Tags))
		for i, tag := range forecast.Tags {
			tags[i] = model.Tag{
				Name: tag.Name,
			}
		}

		// Generate UUIDs and set ForecastID for associated Sources
		sources := make([]model.Source, len(forecast.Sources))
		for k, src := range forecast.Sources {
			sources[k] = model.Source{
				ID:         uuid.New(), // New UUID for each source
				ForecastID: forecastID,
				Name:       src.Name,
				Title:      src.Title,
				URL:        src.URL,
				ImageURL:   &src.ImageURL,
			}
		}

		var market *model.Market
		if forecast.Market != nil {
			market = &model.Market{
				ID:            forecast.Market.ID,
				Question:      forecast.Market.Question,
				Outcomes:      forecast.Market.Outcomes,
				OutcomePrices: forecast.Market.OutcomePrices,
				Volume:        forecast.Market.Volume,
				ImageURL:      forecast.Market.ImageURL,
			}
		}

		// Construct the model forecast with the generated UUID and associations
		modelForecasts[i] = model.Forecast{
			ID:        forecastID, // Set the generated UUID for the forecast
			Headline:  forecast.Headline,
			Summary:   forecast.Summary,
			ImageURL:  forecast.ImageURL,
			Category:  forecast.Category,
			Timestamp: forecast.Timestamp,
			Outcomes:  outcomes,
			Tags:      tags,
			Sources:   sources,
			Market:    market,
		}
	}

	return modelForecasts
}

func ParseOutcomesAndPrices(m *dto.Market) error {
	if err := json.Unmarshal([]byte(m.Outcomes), &m.OutcomeList); err != nil {
		return fmt.Errorf("unable to parse outcomes: %w", err)
	}

	if err := json.Unmarshal([]byte(m.OutcomePrices), &m.OutcomePricesList); err != nil {
		return fmt.Errorf("unable to parse outcomePrices: %w", err)
	}

	return nil
}
