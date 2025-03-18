package promptgen

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
	"text/template"
)

type TemplateType string

const (
	GenerateNewsForecast   TemplateType = "generate_news_forecast"
	GenerateMarketForecast TemplateType = "generate_market_forecast"
	SelectArticles         TemplateType = "select_articles"
	SelectMarkets          TemplateType = "select_markets"
	SelectArticleForEvent  TemplateType = "select_article_for_event"
	ExtractKeywords        TemplateType = "extract_keywords"
)

type PromptTemplate struct {
	templates map[TemplateType]*template.Template
}

func (p *PromptTemplate) CreatePrompt(templateType TemplateType, data interface{}) (string, error) {
	return p.executeTemplate(templateType, data)
}

func (p *PromptTemplate) executeTemplate(templateType TemplateType, data interface{}) (string, error) {
	tmpl, exists := p.templates[templateType]
	if !exists {
		return "", fmt.Errorf("invalid template type: %s", templateType)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("error executing template %q: %v", templateType, err)
	}

	return result.String(), nil
}

func LoadPromptTemplate(filepath string) (*PromptTemplate, error) {
	var rawPrompts struct {
		GenerateNewsForecast   string `yaml:"generate_news_forecast"`
		GenerateMarketForecast string `yaml:"generate_market_forecast"`
		SelectArticles         string `yaml:"select_articles"`
		SelectMarkets          string `yaml:"select_markets"`
		SelectArticleForEvent  string `yaml:"select_article_for_event"`
		ExtractKeywords        string `yaml:"extract_keywords"`
	}

	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(file, &rawPrompts); err != nil {
		return nil, err
	}

	templates := map[TemplateType]*template.Template{}
	for key, value := range map[TemplateType]string{
		GenerateNewsForecast:   rawPrompts.GenerateNewsForecast,
		GenerateMarketForecast: rawPrompts.GenerateMarketForecast,
		SelectArticles:         rawPrompts.SelectArticles,
		SelectMarkets:          rawPrompts.SelectMarkets,
		SelectArticleForEvent:  rawPrompts.SelectArticleForEvent,
		ExtractKeywords:        rawPrompts.ExtractKeywords,
	} {
		tmpl, err := template.New(string(key)).Parse(value)
		if err != nil {
			return nil, fmt.Errorf("error parsing template %q: %v", key, err)
		}
		templates[key] = tmpl
	}

	return &PromptTemplate{templates: templates}, nil
}
