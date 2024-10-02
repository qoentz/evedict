package replicate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
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

func (r *Service) GetPredictions(prompt string, articles []newsapi.Article) (*llm.Predictions, error) {
	if len(prompt) == 0 {
		return nil, fmt.Errorf("empty prompt provided")
	}

	maxTokens := 1024
	payload := RequestPayload{
		Stream: false,
		Input: Input{
			Prompt:    prompt,
			MaxTokens: maxTokens,
		},
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	// Start the prediction
	req, err := http.NewRequest("POST", r.ModelURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.APIKey)

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d\nResponse Body:\n%s", resp.StatusCode, string(body))
	}

	var prediction ResponsePayload
	err = json.NewDecoder(resp.Body).Decode(&prediction)
	if err != nil {
		return nil, fmt.Errorf("error parsing response JSON: %v", err)
	}

	// Poll for prediction completion using the "get" URL from the response
	for prediction.Status != "succeeded" && prediction.Status != "failed" {
		time.Sleep(2 * time.Second)

		// Get the prediction status
		req, err = http.NewRequest("GET", prediction.URLs.Get, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating GET request: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+r.APIKey)

		resp, err = r.HTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making GET request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("unexpected status code: %d\nResponse Body:\n%s", resp.StatusCode, string(body))
		}

		err = json.NewDecoder(resp.Body).Decode(&prediction)
		if err != nil {
			return nil, fmt.Errorf("error parsing prediction JSON: %v", err)
		}
	}

	// Handle the case where prediction.Output is an array of strings
	var outputStr string
	switch v := prediction.Output.(type) {
	case string:
		outputStr = v
	case []interface{}:
		// Concatenate the strings in the array
		var builder strings.Builder
		for _, item := range v {
			if str, ok := item.(string); ok {
				builder.WriteString(str)
			} else {
				return nil, fmt.Errorf("prediction output array contains non-string elements")
			}
		}
		outputStr = builder.String()
	default:
		return nil, fmt.Errorf("unexpected type for prediction output: %T", prediction.Output)
	}

	// Clean up the output string (optional)
	outputStr = strings.TrimSpace(outputStr)

	// Parse the output JSON into Predictions struct
	var predictions llm.Predictions
	err = json.Unmarshal([]byte(outputStr), &predictions)
	if err != nil {
		return nil, fmt.Errorf("error parsing prediction output: %v\nOutput Data:\n%s", err, outputStr)
	}

	for i := range predictions.Predictions {
		predictions.Predictions[i].ImageURL = articles[i].URLToImage
	}

	return &predictions, nil
}

//func InitiateStream(prompt string) (string, error) {
//	if len(prompt) == 0 {
//		return "", fmt.Errorf("empty prompt provided")
//	}
//
//	maxTokens := 1024
//	payload := RequestPayload{
//		Stream: false,
//		Input: Input{
//			Prompt:    prompt,
//			MaxTokens: maxTokens,
//		},
//	}
//
//	reqBody, err := json.Marshal(payload)
//	if err != nil {
//		return "", fmt.Errorf("error marshaling request body: %v", err)
//	}
//
//	resp, err := httputil.PostRequest(os.Getenv("REPLICA_MODEL"), reqBody, os.Getenv("REPLICA_KEY"))
//	if err != nil {
//		return "", err
//	}
//	defer resp.Close()
//
//	body, _ := io.ReadAll(resp)
//	fmt.Println("HERE: ", string(body))
//
//	var response ResponsePayload
//	err = json.NewDecoder(resp).Decode(&response)
//	if err != nil {
//		return "", fmt.Errorf("error decoding API response: %v", err)
//	}
//
//	if response.URLs.Stream == "" {
//		return "", fmt.Errorf("no stream URL found in the response")
//	}
//
//	return response.URLs.Stream, nil
//}
//
//func HandleStream(streamURL string) (*Predictions, error) {
//	req, err := http.NewRequest("GET", streamURL, nil)
//	if err != nil {
//		return nil, fmt.Errorf("error creating request: %v", err)
//	}
//
//	req.Header.Set("Accept", "text/event-stream")
//	req.Header.Set("Cache-Control", "no-store")
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return nil, fmt.Errorf("error connecting to stream: %v", err)
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
//	}
//
//	reader := bufio.NewReader(resp.Body)
//	var output strings.Builder
//
//	for {
//		line, err := reader.ReadString('\n')
//		if err != nil {
//			if err.Error() == "EOF" {
//				break
//			}
//			return nil, fmt.Errorf("error reading stream: %v", err)
//		}
//
//		// Accumulate and simulate writing out each "data:" line in real time
//		line = strings.TrimSpace(line)
//		if strings.HasPrefix(line, "data:") {
//			data := strings.TrimPrefix(line, "data:")
//			output.WriteString(data)
//		} else if strings.HasPrefix(line, "event: done") {
//			break
//		}
//	}
//
//	// Log the raw accumulated JSON data
//	jsonString := output.String()
//	fmt.Println("Accumulated JSON data:", jsonString)
//
//	// Use regex to clean up the keys (remove extra spaces around key names)
//	re := regexp.MustCompile(`"\s*(\w+)\s*"\s*:`)
//	jsonString = re.ReplaceAllString(jsonString, `"$1":`)
//
//	// Convert the cleaned JSON string to a struct
//	var predictions Predictions
//	err = json.Unmarshal([]byte(jsonString), &predictions)
//	if err != nil {
//		log.Printf("Error unmarshaling JSON: %v", err)
//		log.Printf("Cleaned JSON: %s", jsonString) // Output the cleaned JSON for troubleshooting
//		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
//	}
//
//	// Post-process and clean up extra spaces in fields
//	// Regex to remove extra spaces before punctuation
//	spaceBeforePunct := regexp.MustCompile(`\s+([.,])`)
//	multipleSpaces := regexp.MustCompile(`\s{2,}`)
//
//	for i := range predictions.Predictions {
//		// Remove extra spaces between words and before punctuation in title and content
//		predictions.Predictions[i].Title = spaceBeforePunct.ReplaceAllString(predictions.Predictions[i].Title, "$1")
//		predictions.Predictions[i].Title = multipleSpaces.ReplaceAllString(predictions.Predictions[i].Title, " ")
//		predictions.Predictions[i].Content = spaceBeforePunct.ReplaceAllString(predictions.Predictions[i].Content, "$1")
//		predictions.Predictions[i].Content = multipleSpaces.ReplaceAllString(predictions.Predictions[i].Content, " ")
//	}
//
//	// Marshal it back to a clean and formatted JSON string
//	//cleanedJSON, err := json.MarshalIndent(predictions, "", "    ")
//	//if err != nil {
//	//	log.Printf("Error marshaling cleaned JSON: %v", err)
//	//	return nil, fmt.Errorf("error marshaling cleaned JSON: %v", err)
//	//}
//
//	// Output the polished JSON
//	//fmt.Println(string(cleanedJSON))
//
//	return &predictions, nil
//}
