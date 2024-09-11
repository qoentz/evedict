package hugface

import (
	"cashnew/httputil"
	"encoding/json"
	"fmt"
	"os"
)

type RequestPayload struct {
	Inputs     string                 `json:"inputs"`
	Parameters map[string]interface{} `json:"parameters"`
}

type ResponsePayload struct {
	SummaryText string `json:"summary_text"`
}

func GenerateSummary(prompt string) (string, error) {
	if len(prompt) == 0 {
		return "", fmt.Errorf("empty prompt provided")
	}

	payload := RequestPayload{
		Inputs: prompt,
		Parameters: map[string]interface{}{
			"max_length":  1024,
			"min_length":  50,
			"do_sample":   true,
			"temperature": 0.7, // Controls randomness, lower = more focused
			"top_p":       0.9,
		},
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := httputil.PostRequest(os.Getenv("HUGGING_FACE_MODEL"), reqBody, os.Getenv("HUGGING_FACE_KEY"))
	if err != nil {
		return "", err
	}

	fmt.Println(string(resp))

	var responsePayload []ResponsePayload
	if err = json.Unmarshal(resp, &responsePayload); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	return responsePayload[0].SummaryText, nil
}
