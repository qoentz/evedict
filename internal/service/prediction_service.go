package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/qoentz/evedict/internal/db/model"
	"github.com/qoentz/evedict/internal/db/repository"
	"github.com/qoentz/evedict/internal/dto"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"github.com/qoentz/evedict/internal/llm"
	"github.com/qoentz/evedict/internal/llm/replicate"
	"github.com/qoentz/evedict/internal/promptgen"
	"time"
)

type PredictionService struct {
	PredictionRepository *repository.PredictionRepository
	AIService            llm.Service
	NewsAPIService       *newsapi.Service
	PromptTemplate       *promptgen.PromptTemplate
}

func NewPredictionService(predictionRepository *repository.PredictionRepository, replicateService *replicate.Service, newsAPIService *newsapi.Service, template *promptgen.PromptTemplate) *PredictionService {
	return &PredictionService{
		PredictionRepository: predictionRepository,
		AIService:            replicateService,
		NewsAPIService:       newsAPIService,
		PromptTemplate:       template,
	}
}

func (s *PredictionService) GeneratePredictions(category newsapi.Category) ([]dto.Prediction, error) {
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

	var predictions []dto.Prediction
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

		predictionPrompt, err := s.PromptTemplate.CreatePredictionPrompt(mainArticle, articles)
		if err != nil {
			return nil, fmt.Errorf("error creating prediction prompt: %v", err)
		}

		prediction, err := s.AIService.GetPrediction(predictionPrompt)
		if err != nil {
			return nil, fmt.Errorf("error generating prediction: %v", err)
		}

		prediction.ImageURL = mainArticle.URLToImage

		var sources []dto.Source

		mainSource := dto.Source{
			Name:  mainArticle.Source.Name,
			Title: mainArticle.Title,
			URL:   mainArticle.URL,
		}
		sources = append(sources, mainSource)

		for _, article := range articles {
			if article.URL == mainArticle.URL {
				continue
			}

			source := dto.Source{
				Name:  article.Source.Name,
				Title: article.Title,
				URL:   article.URL,
			}

			sources = append(sources, source)
		}

		prediction.Sources = sources
		prediction.Timestamp = time.Now()
		predictions = append(predictions, *prediction)
	}

	return predictions, nil
}

func (s *PredictionService) GetPredictions() ([]dto.Prediction, error) {
	predictions, err := s.PredictionRepository.GetPredictions()
	if err != nil {
		return nil, fmt.Errorf("failed to get predictions: %v", err)
	}

	var result []dto.Prediction
	for _, prediction := range predictions {
		dtoPrediction := s.convertToDTO(prediction)
		result = append(result, dtoPrediction)
	}

	return result, nil
}

func (s *PredictionService) SavePredictions(predictions []dto.Prediction) error {
	modelPredictions := s.convertToModel(predictions)
	err := s.PredictionRepository.SavePredictions(modelPredictions)
	if err != nil {
		return fmt.Errorf("failed to save predictions: %v", err)
	}
	return nil
}

func (s *PredictionService) SavePrediction(prediction *dto.Prediction) error {
	modelPredictions := s.convertToModel([]dto.Prediction{*prediction})
	err := s.PredictionRepository.SavePrediction(&modelPredictions[0])
	if err != nil {
		return fmt.Errorf("failed to save prediction: %v", err)
	}
	return nil
}

func (s *PredictionService) convertToDTO(prediction model.Prediction) dto.Prediction {
	dtoOutcomes := make([]dto.Outcome, len(prediction.Outcomes))
	for i, o := range prediction.Outcomes {
		dtoOutcomes[i] = dto.Outcome{
			Content:         o.Content,
			ConfidenceLevel: o.ConfidenceLevel,
		}
	}

	dtoSources := make([]dto.Source, len(prediction.Sources))
	for i, src := range prediction.Sources {
		dtoSources[i] = dto.Source{
			Name:  src.Name,
			Title: src.Title,
			URL:   src.URL,
		}
	}

	return dto.Prediction{
		Headline:  prediction.Headline,
		Summary:   prediction.Summary,
		Outcomes:  dtoOutcomes,
		ImageURL:  prediction.ImageURL,
		Sources:   dtoSources,
		Timestamp: prediction.Timestamp,
	}
}

func (s *PredictionService) convertToModel(predictions []dto.Prediction) []model.Prediction {
	modelPredictions := make([]model.Prediction, len(predictions))

	for i, prediction := range predictions {
		// Generate a UUID for the prediction
		predictionID := uuid.New()

		// Generate UUIDs and set PredictionID for associated Outcomes
		outcomes := make([]model.Outcome, len(prediction.Outcomes))
		for j, o := range prediction.Outcomes {
			outcomes[j] = model.Outcome{
				ID:              uuid.New(), // New UUID for each outcome
				PredictionID:    predictionID,
				Content:         o.Content,
				ConfidenceLevel: o.ConfidenceLevel,
			}
		}

		// Generate UUIDs and set PredictionID for associated Sources
		sources := make([]model.Source, len(prediction.Sources))
		for k, src := range prediction.Sources {
			sources[k] = model.Source{
				ID:           uuid.New(), // New UUID for each source
				PredictionID: predictionID,
				Name:         src.Name,
				Title:        src.Title,
				URL:          src.URL,
			}
		}

		// Construct the model Prediction with the generated UUID and associations
		modelPredictions[i] = model.Prediction{
			ID:        predictionID, // Set the generated UUID for the prediction
			Headline:  prediction.Headline,
			Summary:   prediction.Summary,
			ImageURL:  prediction.ImageURL,
			Timestamp: prediction.Timestamp,
			Outcomes:  outcomes,
			Sources:   sources,
		}
	}

	return modelPredictions
}
