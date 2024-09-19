package replicate

import (
	"bufio"
	"encoding/json"
	"evedict/internal/httputil"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Prediction struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Predictions struct {
	Predictions []Prediction `json:"predictions"`
}

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

	//body, _ := io.ReadAll(resp)
	//fmt.Println(string(body))

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

func HandleStream(streamURL string) (*Predictions, error) {
	req, err := http.NewRequest("GET", streamURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-store")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error connecting to stream: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)
	var output strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("error reading stream: %v", err)
		}

		// Accumulate and simulate writing out each "data:" line in real time
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimPrefix(line, "data:")
			output.WriteString(data)
		} else if strings.HasPrefix(line, "event: done") {
			break
		}
	}

	// Log the raw accumulated JSON data
	jsonString := output.String()
	fmt.Println("Accumulated JSON data:", jsonString)

	// Use regex to clean up the keys (remove extra spaces around key names)
	re := regexp.MustCompile(`"\s*(\w+)\s*"\s*:`)
	jsonString = re.ReplaceAllString(jsonString, `"$1":`)

	// Convert the cleaned JSON string to a struct
	var predictions Predictions
	err = json.Unmarshal([]byte(jsonString), &predictions)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		log.Printf("Cleaned JSON: %s", jsonString) // Output the cleaned JSON for troubleshooting
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// Post-process and clean up extra spaces in fields
	// Regex to remove extra spaces before punctuation
	spaceBeforePunct := regexp.MustCompile(`\s+([.,])`)
	multipleSpaces := regexp.MustCompile(`\s{2,}`)

	for i := range predictions.Predictions {
		// Remove extra spaces between words and before punctuation in title and content
		predictions.Predictions[i].Title = spaceBeforePunct.ReplaceAllString(predictions.Predictions[i].Title, "$1")
		predictions.Predictions[i].Title = multipleSpaces.ReplaceAllString(predictions.Predictions[i].Title, " ")
		predictions.Predictions[i].Content = spaceBeforePunct.ReplaceAllString(predictions.Predictions[i].Content, "$1")
		predictions.Predictions[i].Content = multipleSpaces.ReplaceAllString(predictions.Predictions[i].Content, " ")
	}

	// Marshal it back to a clean and formatted JSON string
	//cleanedJSON, err := json.MarshalIndent(predictions, "", "    ")
	//if err != nil {
	//	log.Printf("Error marshaling cleaned JSON: %v", err)
	//	return nil, fmt.Errorf("error marshaling cleaned JSON: %v", err)
	//}

	// Output the polished JSON
	//fmt.Println(string(cleanedJSON))

	return &predictions, nil
}
