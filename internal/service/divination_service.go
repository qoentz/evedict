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

type DivinationService struct {
	DivinationRepository *repository.DivinationRepository
	AIService            llm.Service
	NewsAPIService       *newsapi.Service
	PromptTemplate       *promptgen.PromptTemplate
}

func NewDivinationService(divinationRepository *repository.DivinationRepository, replicateService *replicate.Service, newsAPIService *newsapi.Service, template *promptgen.PromptTemplate) *DivinationService {
	return &DivinationService{
		DivinationRepository: divinationRepository,
		AIService:            replicateService,
		NewsAPIService:       newsAPIService,
		PromptTemplate:       template,
	}
}

func (s *DivinationService) GenerateDivinations(category newsapi.Category) ([]dto.Divination, error) {
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

	var divinations []dto.Divination
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

		divinationPrompt, err := s.PromptTemplate.CreateDivinationPrompt(mainArticle, articles)
		if err != nil {
			return nil, fmt.Errorf("error creating divination prompt: %v", err)
		}

		divination, err := s.AIService.GetDivination(divinationPrompt)
		if err != nil {
			return nil, fmt.Errorf("error generating divination: %v", err)
		}

		divination.ImageURL = mainArticle.URLToImage

		var sources []dto.Source

		mainSource := dto.Source{
			Name:  mainArticle.Source.Name,
			Title: mainArticle.Title,
			URL:   mainArticle.URL,
		}
		sources = append(sources, mainSource)

		for _, article := range articles {
			if article.URL == mainArticle.URL || article.Title == "[Removed]" {
				continue
			}

			source := dto.Source{
				Name:  article.Source.Name,
				Title: article.Title,
				URL:   article.URL,
			}

			sources = append(sources, source)
		}

		divination.Sources = sources
		divination.Timestamp = time.Now().UTC()
		divinations = append(divinations, *divination)
	}

	return divinations, nil
}

func (s *DivinationService) GetDivinations() ([]dto.Divination, error) {
	divinations, err := s.DivinationRepository.GetDivinations()
	if err != nil {
		return nil, fmt.Errorf("failed to get divinations: %v", err)
	}

	var result []dto.Divination
	for _, divination := range divinations {
		dtoDivination := s.convertToDTO(divination)
		result = append(result, dtoDivination)
	}

	return result, nil
}

func (s *DivinationService) SaveDivinations(divinations []dto.Divination) error {
	modelDivinations := s.convertToModel(divinations)
	err := s.DivinationRepository.SaveDivinations(modelDivinations)
	if err != nil {
		return fmt.Errorf("failed to save divinations: %v", err)
	}
	return nil
}

func (s *DivinationService) SaveDivination(divination *dto.Divination) error {
	modelDivinations := s.convertToModel([]dto.Divination{*divination})
	err := s.DivinationRepository.SaveDivination(&modelDivinations[0])
	if err != nil {
		return fmt.Errorf("failed to save divination: %v", err)
	}
	return nil
}

func (s *DivinationService) convertToDTO(divination model.Divination) dto.Divination {
	dtoOutcomes := make([]dto.Outcome, len(divination.Outcomes))
	for i, o := range divination.Outcomes {
		dtoOutcomes[i] = dto.Outcome{
			Content:         o.Content,
			ConfidenceLevel: o.ConfidenceLevel,
		}
	}

	dtoSources := make([]dto.Source, len(divination.Sources))
	for i, src := range divination.Sources {
		dtoSources[i] = dto.Source{
			Name:  src.Name,
			Title: src.Title,
			URL:   src.URL,
		}
	}

	return dto.Divination{
		Headline:  divination.Headline,
		Summary:   divination.Summary,
		Outcomes:  dtoOutcomes,
		ImageURL:  divination.ImageURL,
		Sources:   dtoSources,
		Timestamp: divination.Timestamp,
	}
}

func (s *DivinationService) convertToModel(divinations []dto.Divination) []model.Divination {
	modelDivinations := make([]model.Divination, len(divinations))

	for i, divination := range divinations {
		// Generate a UUID for the divination
		divinationID := uuid.New()

		// Generate UUIDs and set DivinationID for associated Outcomes
		outcomes := make([]model.Outcome, len(divination.Outcomes))
		for j, o := range divination.Outcomes {
			outcomes[j] = model.Outcome{
				ID:              uuid.New(), // New UUID for each outcome
				DivinationID:    divinationID,
				Content:         o.Content,
				ConfidenceLevel: o.ConfidenceLevel,
			}
		}

		// Generate UUIDs and set DivinationID for associated Sources
		sources := make([]model.Source, len(divination.Sources))
		for k, src := range divination.Sources {
			sources[k] = model.Source{
				ID:           uuid.New(), // New UUID for each source
				DivinationID: divinationID,
				Name:         src.Name,
				Title:        src.Title,
				URL:          src.URL,
			}
		}

		// Construct the model Divination with the generated UUID and associations
		modelDivinations[i] = model.Divination{
			ID:        divinationID, // Set the generated UUID for the divination
			Headline:  divination.Headline,
			Summary:   divination.Summary,
			ImageURL:  divination.ImageURL,
			Timestamp: divination.Timestamp,
			Outcomes:  outcomes,
			Sources:   sources,
		}
	}

	return modelDivinations
}
