package replicate

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-wordwrap"
	"io"
	"net/http"
	"os"
	"strings"
)

type FinalPredictionResponse struct {
	Status string   `json:"status"`
	Output []string `json:"output"`
	Error  string   `json:"error"`
}

// CheckPredictionStatus makes a single GET request to check the prediction status
func CheckPredictionStatus(getURL string) (string, error) {
	// Make the GET request to check the status of the prediction
	req, err := http.NewRequest("GET", getURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set the Authorization header with the API key
	req.Header.Set("Authorization", "Bearer "+os.Getenv("REPLICA_KEY"))

	// Make the GET request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	var finalResp FinalPredictionResponse
	if err := json.Unmarshal(body, &finalResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	// Check if the status is "succeeded" and return the combined output
	if finalResp.Status == "succeeded" {
		return strings.Join(finalResp.Output, ""), nil
	} else if finalResp.Status == "failed" {
		return "", fmt.Errorf("prediction failed: %s", finalResp.Error)
	} else {
		return fmt.Sprintf("Status: %s", finalResp.Status), nil
	}
}

func FormatPrediction(prediction string) string {
	// Define ANSI escape codes for bold text
	boldStart := "\033[1m"
	boldEnd := "\033[0m"

	// Break the content into sections
	sections := strings.Split(prediction, "**")

	// Format sections
	var formatted string
	for _, section := range sections {
		if strings.HasPrefix(section, " ") {
			// Headings (e.g., "Gaming and AI", "Crypto and Web3")
			formatted += boldStart + wordwrap.WrapString(strings.TrimSpace(section), 80) + boldEnd + "\n\n"
		} else {
			// Regular text, apply word wrap
			formatted += wordwrap.WrapString(strings.TrimSpace(section), 80) + "\n\n"
		}
	}

	return formatted
}
