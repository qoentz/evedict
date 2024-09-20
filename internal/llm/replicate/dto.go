package replicate

type RequestPayload struct {
	Stream bool  `json:"stream"`
	Input  Input `json:"input"`
}

type Input struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type ResponsePayload struct {
	ID      string      `json:"id"`
	Model   string      `json:"model"`
	Version string      `json:"version"`
	Status  string      `json:"status"`
	Output  interface{} `json:"output"`
	URLs    URLs        `json:"urls"`
	Stream  bool        `json:"stream"`
}

type URLs struct {
	Cancel string `json:"cancel"`
	Get    string `json:"get"`
	Stream string `json:"stream"`
}

type Predictions struct {
	Predictions []Prediction `json:"predictions"`
}

type Prediction struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
