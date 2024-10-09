package llm

import (
	"github.com/qoentz/evedict/internal/dto"
)

type Service interface {
	GetPrediction(prompt string) (*dto.Prediction, error)
	SelectArticles(prompt string) ([]int, error)
	ExtractKeywords(prompt string) ([]string, error)
}
