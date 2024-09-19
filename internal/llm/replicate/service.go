package replicate

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
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
	ID      string      `json:"id"`
	Model   string      `json:"model"`
	Version string      `json:"version"`
	Status  string      `json:"status"`
	Output  interface{} `json:"output"`
	URLs    URLs        `json:"urls"`   // Correctly map the "urls" object
	Stream  bool        `json:"stream"` // Handle the top-level "stream" field as a bool
}

// Input represents the input field of the JSON request payload
type Input struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
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

func GetPredictions(prompt string) (*Predictions, error) {
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
	req, err := http.NewRequest("POST", os.Getenv("REPLICA_MODEL"), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("REPLICA_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
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
		req, err := http.NewRequest("GET", prediction.URLs.Get, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating GET request: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+os.Getenv("REPLICA_KEY"))

		resp, err := client.Do(req)
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
	var predictions Predictions
	err = json.Unmarshal([]byte(outputStr), &predictions)
	if err != nil {
		return nil, fmt.Errorf("error parsing prediction output: %v\nOutput Data:\n%s", err, outputStr)
	}

	return &predictions, nil
}

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
//	var predictions Predictions
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
//		line = strings.TrimSpace(line)
//		if strings.HasPrefix(line, "data:") {
//			data := strings.TrimPrefix(line, "data:")
//			data = strings.TrimSpace(data)
//
//			// Skip empty data lines
//			if data == "" || data == "[DONE]" {
//				continue
//			}
//
//			// Parse the JSON data
//			var pred Prediction
//			err := json.Unmarshal([]byte(data), &pred)
//			if err != nil {
//				return nil, fmt.Errorf("error parsing JSON: %v", err)
//			}
//
//			// Append the prediction to the slice
//			predictions.Predictions = append(predictions.Predictions, pred)
//		} else if strings.HasPrefix(line, "event: done") {
//			break
//		}
//	}
//
//	// Return the accumulated Predictions
//	return &predictions, nil
//}

//func HandleStream(streamURL string) (*Predictions, error) {
//	// Create the request
//	req, err := http.NewRequest("GET", streamURL, nil)
//	if err != nil {
//		return nil, fmt.Errorf("error creating request: %v", err)
//	}
//
//	// Set headers for SSE
//	req.Header.Set("Accept", "text/event-stream")
//	req.Header.Set("Cache-Control", "no-store")
//
//	// Execute the request
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return nil, fmt.Errorf("error connecting to stream: %v", err)
//	}
//	defer resp.Body.Close()
//
//	// Check for HTTP errors
//	if resp.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
//	}
//
//	// Read the stream
//	reader := bufio.NewReader(resp.Body)
//	var buffer strings.Builder
//
//	for {
//		line, err := reader.ReadString('\n')
//		if err != nil {
//			if err == io.EOF {
//				break
//			}
//			return nil, fmt.Errorf("error reading stream: %v", err)
//		}
//
//		// Remove any trailing carriage return or newline characters
//		line = strings.TrimRight(line, "\r\n")
//
//		if strings.HasPrefix(line, "data:") {
//			// Extract the data after "data:"
//			dataLine := strings.TrimPrefix(line, "data:")
//			// Append dataLine to buffer without adding spaces
//			buffer.WriteString(dataLine)
//		} else if strings.HasPrefix(line, "event: done") {
//			// End of stream
//			break
//		}
//	}
//
//	// Get the accumulated JSON data
//	jsonData := buffer.String()
//	jsonData = strings.TrimSpace(jsonData)
//
//	// Remove any newline and carriage return characters
//	jsonData = strings.ReplaceAll(jsonData, "\n", "")
//	jsonData = strings.ReplaceAll(jsonData, "\r", "")
//
//	// Print the accumulated JSON data (optional, for debugging)
//	fmt.Println("Accumulated JSON Data:")
//	fmt.Println(jsonData)
//
//	// Parse the JSON data into Predictions struct
//	var predictions Predictions
//	err = json.Unmarshal([]byte(jsonData), &predictions)
//	if err != nil {
//		return nil, fmt.Errorf("error parsing JSON: %v\nJSON Data:\n%s", err, jsonData)
//	}
//
//	return &predictions, nil
//}

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
