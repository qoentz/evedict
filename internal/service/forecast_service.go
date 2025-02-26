package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/qoentz/evedict/internal/api/dto"
	"github.com/qoentz/evedict/internal/db/model"
	"github.com/qoentz/evedict/internal/db/repository"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"github.com/qoentz/evedict/internal/llm"
	"github.com/qoentz/evedict/internal/llm/replicate"
	"github.com/qoentz/evedict/internal/promptgen"
	"time"
)

type ForecastService struct {
	ForecastRepository *repository.ForecastRepository
	AIService          llm.Service
	NewsAPIService     *newsapi.Service
	PromptTemplate     *promptgen.PromptTemplate
}

func NewForecastService(forecastRepository *repository.ForecastRepository, replicateService *replicate.Service, newsAPIService *newsapi.Service, template *promptgen.PromptTemplate) *ForecastService {
	return &ForecastService{
		ForecastRepository: forecastRepository,
		AIService:          replicateService,
		NewsAPIService:     newsAPIService,
		PromptTemplate:     template,
	}
}

func (s *ForecastService) GenerateForecasts(category newsapi.Category) ([]dto.Forecast, error) {
	headlines, err := s.NewsAPIService.FetchTopHeadlines(category)
	if err != nil {
		return nil, fmt.Errorf("error fetching headlines from NewsAPI: %v", err)
	}

	selectionPrompt, err := s.PromptTemplate.CreateArticleSelectionPrompt(headlines)
	if err != nil {
		return nil, fmt.Errorf("error creating article selection prompt: %v", err)
	}

	articleSelection, err := s.AIService.SelectArticles(selectionPrompt)
	if err != nil {
		return nil, fmt.Errorf("error selecting articles: %v", err)
	}

	var forecasts []dto.Forecast
	for _, idx := range articleSelection {
		mainArticle := headlines[idx]

		extractionPrompt, err := s.PromptTemplate.CreateKeywordExtractionPrompt(mainArticle)
		if err != nil {
			return nil, fmt.Errorf("error creating keyword extraction prompt: %v", err)
		}

		keywords, err := s.AIService.ExtractKeywords(extractionPrompt)
		if err != nil {
			return nil, fmt.Errorf("error extracting keywords: %v", err)
		}

		articles, err := s.NewsAPIService.FetchWithKeywords(keywords)
		if err != nil {
			return nil, fmt.Errorf("error fetching articles from NewsAPI with keywords: %v", err)
		}

		forecastPrompt, err := s.PromptTemplate.CreateForecastPrompt(mainArticle, articles)
		if err != nil {
			return nil, fmt.Errorf("error creating forecast prompt: %v", err)
		}

		forecast, err := s.AIService.GetForecast(forecastPrompt)
		if err != nil {
			return nil, fmt.Errorf("error generating forecast: %v", err)
		}

		forecast.ImageURL = mainArticle.URLToImage
		forecast.Tags = keywords

		var sources []dto.Source

		mainSource := dto.Source{
			Name:     mainArticle.Source.Name,
			Title:    mainArticle.Title,
			URL:      mainArticle.URL,
			ImageURL: mainArticle.URLToImage,
		}

		sources = append(sources, mainSource)

		for _, article := range articles {
			if article.URL == mainArticle.URL || article.Title == "[Removed]" {
				continue
			}

			source := dto.Source{
				Name:     article.Source.Name,
				Title:    article.Title,
				URL:      article.URL,
				ImageURL: article.URLToImage,
			}

			sources = append(sources, source)
		}

		forecast.Sources = sources
		forecast.Timestamp = time.Now().UTC()
		forecasts = append(forecasts, *forecast)
	}

	return forecasts, nil
}

func (s *ForecastService) GetForecasts(limit int, offset int) ([]dto.Forecast, error) {
	forecasts, err := s.ForecastRepository.GetForecasts(limit, offset)
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

func (s *ForecastService) GetForecast(forecastId uuid.UUID) (*dto.Forecast, error) {
	forecast, err := s.ForecastRepository.GetForecast(forecastId)
	if err != nil {
		return nil, fmt.Errorf("failed to get forecast: %v", err)
	}

	return s.convertToDTO(forecast), nil
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

	return &dto.Forecast{
		ID:        forecast.ID,
		Headline:  forecast.Headline,
		Summary:   forecast.Summary,
		Outcomes:  dtoOutcomes,
		ImageURL:  forecast.ImageURL,
		Sources:   dtoSources,
		Timestamp: forecast.Timestamp,
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
				Name: tag,
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
		}
	}

	return modelForecasts
}
