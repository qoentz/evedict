package llm

import (
	"github.com/qoentz/evedict/internal/api/dto"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"github.com/qoentz/evedict/internal/eventfeed/polymarket"
	"github.com/qoentz/evedict/internal/promptgen"
)

type Service interface {
	GetForecast(mainArticle newsapi.Article, relatedArticles []newsapi.Article, event *polymarket.Event) (*dto.Forecast, error)
	SelectIndexes(templateType promptgen.TemplateType, data interface{}, minSelection int) ([]int, error)
	SelectIndex(templateType promptgen.TemplateType, data interface{}) (int, error)
	ExtractKeywords(article newsapi.Article) ([]string, error)
}
