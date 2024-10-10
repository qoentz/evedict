package hugface

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func TranslateHeadline(headline string) (string, error) {
	reqBody := []byte(`{"inputs": "` + headline + `"}`)

	req, err := http.NewRequest("POST", os.Getenv("TRANSLATE_MODEL"), bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("HUGGING_FACE_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result []map[string]string
	if err = json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	return result[0]["translation_text"], nil
}
