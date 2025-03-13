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
	PromptTemplate     *promptgen.PromptTemplate
	PolyMarketService  *polymarket.Service
}

func NewForecastService(forecastRepository *repository.ForecastRepository, replicateService *replicate.Service, newsAPIService *newsapi.Service, template *promptgen.PromptTemplate, polyMarketService *polymarket.Service) *ForecastService {
	return &ForecastService{
		ForecastRepository: forecastRepository,
		AIService:          replicateService,
		NewsAPIService:     newsAPIService,
		PromptTemplate:     template,
		PolyMarketService:  polyMarketService,
	}
}

func (s *ForecastService) GeneratePolyForecasts() ([]dto.Forecast, error) {
	// fetch events
	events, err := s.PolyMarketService.FetchTopEvents()
	if err != nil {
		return nil, fmt.Errorf("error fetching events: %v", err)
	}

	var SMPEvents []polymarket.Event
	for _, e := range events {
		if len(e.Markets) == 1 {
			SMPEvents = append(SMPEvents, e)
		}
	}

	// extract a relevant market
	marketSelection, err := s.PromptTemplate.CreateMarketSelectionPrompt(SMPEvents)
	if err != nil {
		return nil, fmt.Errorf("error creating market selection prompt: %v", err)
	}

	eventSelection, err := s.AIService.SelectArticles(marketSelection)
	if err != nil {
		return nil, fmt.Errorf("error selecting articles: %v", err)
	}

	fmt.Println(eventSelection)

	var forecasts []dto.Forecast
	for _, idx := range eventSelection {
		mainEvent := events[idx]

		var keywords []string
		for _, tag := range mainEvent.Tags {
			keywords = append(keywords, tag.Label)
		}

		// with the tags from this market, fetch news (on keywords)
		articles, err := s.NewsAPIService.FetchWithKeywords(keywords)
		if err != nil {
			return nil, fmt.Errorf("error fetching articles from NewsAPI with keywords: %v", err)
		}

		eventArticlePrompt, err := s.PromptTemplate.CreateEventArticlePrompt(events[idx], articles)
		if err != nil {
			return nil, fmt.Errorf("error creating eventArticle prompt: %v", err)
		}

		mainArticleIdx, err := s.AIService.SelectArticle(eventArticlePrompt)
		if err != nil {
			return nil, fmt.Errorf("error selecting main article: %v", err)
		}

		forecastPrompt, err := s.PromptTemplate.CreatePolyForecastPrompt(articles[mainArticleIdx], articles, mainEvent)
		if err != nil {
			return nil, fmt.Errorf("error creating forecast prompt: %v", err)
		}

		forecast, err := s.AIService.GetForecast(forecastPrompt)
		if err != nil {
			return nil, fmt.Errorf("error generating forecast: %v", err)
		}

		var firstMarket polymarket.Market
		if len(mainEvent.Markets) > 0 {
			firstMarket = mainEvent.Markets[0]
		}

		// Market json
		forecast.Market = &dto.Market{
			Question:      firstMarket.Question,
			Outcomes:      firstMarket.Outcomes,
			OutcomePrices: firstMarket.OutcomePrices,
			Volume:        firstMarket.Volume,
			ImageURL:      mainEvent.Image,
			ExternalID:    firstMarket.ID,
		}

		forecast.ImageURL = articles[mainArticleIdx].URLToImage

		tags := make([]dto.Tag, len(keywords))
		for i, t := range keywords {
			tags[i].Name = t
		}

		forecast.Tags = tags

		var sources []dto.Source

		mainSource := dto.Source{
			Name:     articles[mainArticleIdx].Source.Name,
			Title:    articles[mainArticleIdx].Title,
			URL:      articles[mainArticleIdx].URL,
			ImageURL: articles[mainArticleIdx].URLToImage,
		}

		sources = append(sources, mainSource)

		for _, article := range articles {
			if article.URL == articles[mainArticleIdx].URL || article.Title == "[Removed]" {
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

		exists, _ := s.ForecastRepository.CheckImageURL(mainArticle.URLToImage)
		if exists {
			log.Println(mainArticle.Title + " already exists!")
			continue
		}

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

		tags := make([]dto.Tag, len(keywords))
		for i, t := range keywords {
			tags[i].Name = t
		}

		forecast.Tags = tags

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
	fmt.Println("CATEGORY: ", modelForecasts[0].Category, " ", modelForecasts[1].Category)
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
				ID:            uuid.New(),
				Question:      forecast.Market.Question,
				Outcomes:      forecast.Market.Outcomes,
				OutcomePrices: forecast.Market.OutcomePrices,
				Volume:        forecast.Market.Volume,
				ImageURL:      forecast.Market.ImageURL,
				ExternalID:    forecast.Market.ExternalID,
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
