package llm

import (
	"github.com/qoentz/evedict/internal/api/dto"
)

type Service interface {
	GetForecast(prompt string) (*dto.Forecast, error)
	SelectArticles(prompt string) ([]int, error)
	ExtractKeywords(prompt string) ([]string, error)
}
