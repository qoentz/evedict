package llm

type Service interface {
	GetPredictions(prompt string) (*Predictions, error)
}

type Predictions struct {
	Predictions []Prediction `json:"predictions"`
}

type Prediction struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
