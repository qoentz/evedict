package llm

type Service interface {
	GetPrediction(prompt string) (*Prediction, error)
	SelectArticles(prompt string) ([]int, error)
	ExtractKeywords(prompt string) ([]string, error)
}

type Prediction struct {
	Headline string    `json:"headline"`
	Summary  string    `json:"summary"`
	Outcomes []Outcome `json:"outcomes"`
	ImageURL string
}

type Outcome struct {
	Content         string `json:"content"`
	ConfidenceLevel int    `json:"confidenceLevel"`
}
