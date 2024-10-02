package llm

import "github.com/qoentz/evedict/internal/eventfeed/newsapi"

type Service interface {
	GetPredictions(prompt string, articles []newsapi.Article) (*Predictions, error)
}

type Predictions struct {
	Predictions []Prediction `json:"predictions"`
}

type Prediction struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	ImageURL string `json:"imageUrl"`
}
