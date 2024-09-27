package promptgen

import (
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
	"text/template"
)

type PromptTemplate struct {
	PredictEvents string `yaml:"predict_events"`
}

type ArticleData struct {
	Title   string
	Content string
}

func (p *PromptTemplate) CreatePromptWithArticles(articles []newsapi.Article) (string, error) {
	var articleData []ArticleData
	for i, article := range articles {
		if i >= 5 {
			break
		}

		if article.Description != "" {
			articleData = append(articleData, ArticleData{
				Title:   article.Title,
				Content: article.Description,
			})
		}
	}

	tmpl, err := template.New("prompt").Parse(p.PredictEvents)
	if err != nil {
		return "", err
	}
	var result strings.Builder
	err = tmpl.Execute(&result, map[string]interface{}{
		"Articles": articleData,
	})
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func LoadPromptTemplate(filepath string) (*PromptTemplate, error) {
	var prompts PromptTemplate
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &prompts)
	if err != nil {
		return nil, err
	}
	return &prompts, nil
}
