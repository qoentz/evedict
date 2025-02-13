package llm

import (
	"github.com/qoentz/evedict/internal/api/dto"
)

type Service interface {
	GetDivination(prompt string) (*dto.Divination, error)
	SelectArticles(prompt string) ([]int, error)
	ExtractKeywords(prompt string) ([]string, error)
}
