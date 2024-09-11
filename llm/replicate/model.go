package replicate

import (
	"encoding/json"
	"evedict/httputil"
	"fmt"
	"os"
)

type RequestPayload struct {
	Stream bool  `json:"stream"`
	Input  Input `json:"input"`
}

// Input represents the input field of the JSON request payload
type Input struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

func GenerateSummary(prompt string) (string, error) {
	if len(prompt) == 0 {
		return "", fmt.Errorf("empty prompt provided")
	}

	maxTokens := 1024
	payload := RequestPayload{
		Stream: true,
		Input: Input{
			Prompt:    prompt,
			MaxTokens: maxTokens,
		},
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := httputil.PostRequest(os.Getenv("REPLICA_MODEL"), reqBody, os.Getenv("REPLICA_KEY"))
	if err != nil {
		return "", err
	}

	return string(resp), nil
}
