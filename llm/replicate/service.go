package replicate

import (
	"bufio"
	"encoding/json"
	"evedict/httputil"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type RequestPayload struct {
	Stream bool  `json:"stream"`
	Input  Input `json:"input"`
}

type URLs struct {
	Cancel string `json:"cancel"`
	Get    string `json:"get"`
	Stream string `json:"stream"`
}

type ResponsePayload struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Version string `json:"version"`
	URLs    URLs   `json:"urls"`   // Correctly map the "urls" object
	Stream  bool   `json:"stream"` // Handle the top-level "stream" field as a bool
}

// Input represents the input field of the JSON request payload
type Input struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

func InitiateStream(prompt string) (string, error) {
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
	defer resp.Close()

	var response ResponsePayload
	err = json.NewDecoder(resp).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("error decoding API response: %v", err)
	}

	if response.URLs.Stream == "" {
		return "", fmt.Errorf("no stream URL found in the response")
	}

	return response.URLs.Stream, nil
}

func HandleStream(streamURL string) error {
	req, err := http.NewRequest("GET", streamURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-store")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error connecting to stream: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)
	var output strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("error reading stream: %v", err)
		}

		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "data:") {
			data := strings.TrimPrefix(line, "data: ")
			data = strings.ReplaceAll(data, "data:", "") // Remove any "data:" inside the text
			if data == "{}" || data == "" {
				continue // Skip empty or irrelevant data
			}
			// Accumulate text
			output.WriteString(data)
		} else if line == "" && output.Len() > 0 {
			// When we hit an empty line, output the accumulated text
			typeOutText(output.String()) // Print with typing effect
			output.Reset()               // Reset for the next chunk
		}
	}

	return nil
}

func typeOutText(text string) {
	for _, char := range text {
		fmt.Printf("%c", char)
		time.Sleep(25 * time.Millisecond)
	}
}
