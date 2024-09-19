package hugface

import (
	"encoding/json"
	"evedict/internal/httputil"
	"io"
	"os"
)

func TranslateHeadline(headline string) (string, error) {
	reqBody := []byte(`{"inputs": "` + headline + `"}`)

	resp, err := httputil.PostRequest(os.Getenv("TRANSLATE_MODEL"), reqBody, os.Getenv("HUGGING_FACE_KEY"))
	if err != nil {
		return "", err
	}
	defer resp.Close()

	respBody, err := io.ReadAll(resp)
	if err != nil {
		return "", err
	}

	var result []map[string]string
	if err = json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	return result[0]["translation_text"], nil
}
