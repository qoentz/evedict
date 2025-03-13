package promptgen

import (
	"fmt"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"github.com/qoentz/evedict/internal/eventfeed/polymarket"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
	"text/template"
)

type PromptTemplate struct {
	PredictEvents         string `yaml:"predict_events"`
	PredictPolyEvents     string `yaml:"predict_poly_events"`
	SelectArticles        string `yaml:"select_articles"`
	SelectMarkets         string `yaml:"select_markets"`
	SelectArticleForEvent string `yaml:"select_article_for_event"`
	ExtractKeywords       string `yaml:"extract_keywords"`
}

func (p *PromptTemplate) CreatePolyForecastPrompt(mainArticle newsapi.Article, relatedArticles []newsapi.Article, event polymarket.Event) (string, error) {
	if mainArticle.Title == "" || mainArticle.Description == "" {
		return "", fmt.Errorf("main article is missing title or description")
	}

	t, err := template.New("forecastPrompt").Parse(p.PredictEvents)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	data := struct {
		MainArticle     newsapi.Article
		RelatedArticles []newsapi.Article
		Event           polymarket.Event
	}{
		MainArticle:     mainArticle,
		RelatedArticles: relatedArticles,
		Event:           event,
	}

	var promptBuilder strings.Builder
	err = t.Execute(&promptBuilder, data)
	if err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return promptBuilder.String(), nil
}

func (p *PromptTemplate) CreateForecastPrompt(mainArticle newsapi.Article, relatedArticles []newsapi.Article) (string, error) {
	if mainArticle.Title == "" || mainArticle.Description == "" {
		return "", fmt.Errorf("main article is missing title or description")
	}

	t, err := template.New("forecastPrompt").Parse(p.PredictEvents)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	data := struct {
		MainArticle     newsapi.Article
		RelatedArticles []newsapi.Article
	}{
		MainArticle:     mainArticle,
		RelatedArticles: relatedArticles,
	}

	var promptBuilder strings.Builder
	err = t.Execute(&promptBuilder, data)
	if err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return promptBuilder.String(), nil
}

func (p *PromptTemplate) CreateArticleSelectionPrompt(articles []newsapi.Article) (string, error) {
	tmpl, err := template.New("selectArticles").Parse(p.SelectArticles)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	data := struct {
		Articles []newsapi.Article
	}{
		Articles: articles,
	}

	var promptBuilder strings.Builder
	err = tmpl.Execute(&promptBuilder, data)
	if err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return promptBuilder.String(), nil
}

func (p *PromptTemplate) CreateMarketSelectionPrompt(events []polymarket.Event) (string, error) {
	tmpl, err := template.New("selectMarkets").Parse(p.SelectMarkets)
	if err != nil {
		return "", fmt.Errorf("error parsing SelectMarkets template: %v", err)
	}

	data := struct {
		Events []polymarket.Event
	}{
		Events: events,
	}

	var promptBuilder strings.Builder
	err = tmpl.Execute(&promptBuilder, data)
	if err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return promptBuilder.String(), nil
}

func (p *PromptTemplate) CreateEventArticlePrompt(event polymarket.Event, articles []newsapi.Article) (string, error) {
	tmpl, err := template.New("selectArticleForEvent").Parse(p.SelectArticleForEvent)
	if err != nil {
		return "", fmt.Errorf("error parsing selectArticleForEvent template: %v", err)
	}

	data := struct {
		Event    polymarket.Event
		Articles []newsapi.Article
	}{
		Event:    event,
		Articles: articles,
	}

	var promptBuilder strings.Builder
	err = tmpl.Execute(&promptBuilder, data)
	if err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return promptBuilder.String(), nil
}

func (p *PromptTemplate) CreateKeywordExtractionPrompt(article newsapi.Article) (string, error) {
	tmpl, err := template.New("keywordPrompt").Parse(p.ExtractKeywords)
	if err != nil {
		return "", err
	}
	var result strings.Builder
	err = tmpl.Execute(&result, article)
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
