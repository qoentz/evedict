package replicate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/qoentz/evedict/internal/api/dto"
	"github.com/qoentz/evedict/internal/llm"
	"io"
	"net/http"
	"strings"
	"time"
)

type Service struct {
	HTTPClient *http.Client
	ModelURL   string
	APIKey     string
}

var _ llm.Service = &Service{}

func NewReplicateService(client *http.Client, modelURL string, apiKey string) *Service {
	return &Service{
		HTTPClient: client,
		ModelURL:   modelURL,
		APIKey:     apiKey,
	}
}

func (s *Service) GetForecast(prompt string) (*dto.Forecast, error) {
	output, err := s.processRequest(prompt, 1024)
	if err != nil {
		return nil, err
	}

	var result dto.Forecast
	err = json.Unmarshal([]byte(output), &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing forecast output: %v\nOutput Data:\n%s", err, output)
	}

	return &result, nil
}

func (s *Service) SelectArticles(prompt string) ([]int, error) {
	outputStr, err := s.processRequest(prompt, 100)
	if err != nil {
		return nil, err
	}

	type SelectionResponse struct {
		Selected []int `json:"selected"`
	}

	var selection SelectionResponse
	err = json.Unmarshal([]byte(outputStr), &selection)
	if err != nil {
		return nil, fmt.Errorf("error parsing selection output: %v\nOutput Data:\n%s", err, outputStr)
	}

	if len(selection.Selected) < 2 {
		return nil, fmt.Errorf("expected 2 selected articles, got %d: %v", len(selection.Selected), selection.Selected)
	}

	return selection.Selected, nil
}

func (s *Service) ExtractKeywords(prompt string) ([]string, error) {
	outputStr, err := s.processRequest(prompt, 50)
	if err != nil {
		return nil, err
	}

	keywords := strings.Split(strings.TrimSpace(outputStr), ",")
	if len(keywords) != 2 {
		return nil, fmt.Errorf("expected 2 keywords, got %d: %v", len(keywords), keywords)
	}

	for i := range keywords {
		keywords[i] = strings.TrimSpace(keywords[i])
	}

	return keywords, nil
}

func (s *Service) processRequest(prompt string, maxTokens int) (string, error) {
	if len(prompt) == 0 {
		return "", fmt.Errorf("empty prompt provided")
	}

	// Construct the payload for the request
	payload := RequestPayload{
		Stream: false,
		Input: Input{
			Prompt:    prompt,
			MaxTokens: maxTokens,
		},
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %v", err)
	}

	// Make the POST request
	req, err := http.NewRequest("POST", s.ModelURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d\nResponse Body:\n%s", resp.StatusCode, string(body))
	}

	// Decode the initial response to get the forecast status and URLs
	var forecast ResponsePayload
	err = json.NewDecoder(resp.Body).Decode(&forecast)
	if err != nil {
		return "", fmt.Errorf("error parsing response JSON: %v", err)
	}

	// Poll for completion if needed
	for forecast.Status != "succeeded" && forecast.Status != "failed" {
		time.Sleep(2 * time.Second)

		// Poll the status using the "get" URL
		req, err = http.NewRequest("GET", forecast.URLs.Get, nil)
		if err != nil {
			return "", fmt.Errorf("error creating GET request: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+s.APIKey)

		resp, err = s.HTTPClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("error making GET request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return "", fmt.Errorf("unexpected status code: %d\nResponse Body:\n%s", resp.StatusCode, string(body))
		}

		// Update the forecast with the new status
		err = json.NewDecoder(resp.Body).Decode(&forecast)
		if err != nil {
			return "", fmt.Errorf("error parsing forecast JSON: %v", err)
		}
	}

	// Handle the output
	var outputStr string
	switch v := forecast.Output.(type) {
	case string:
		outputStr = v
	case []interface{}:
		// Concatenate the strings in the array
		var builder strings.Builder
		for _, item := range v {
			if str, ok := item.(string); ok {
				builder.WriteString(str)
			} else {
				return "", fmt.Errorf("forecast output array contains non-string elements")
			}
		}
		outputStr = builder.String()
	default:
		return "", fmt.Errorf("unexpected type for forecast output: %T", forecast.Output)
	}

	return strings.TrimSpace(outputStr), nil
}
