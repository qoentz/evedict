package hugface

import (
	"cashnew/httputil"
	"encoding/json"
	"os"
)

func TranslateHeadline(headline string) (string, error) {
	reqBody := []byte(`{"inputs": "` + headline + `"}`)

	resp, err := httputil.PostRequest(os.Getenv("TRANSLATE_MODEL"), reqBody, os.Getenv("HUGGING_FACE_KEY"))
	if err != nil {
		return "", err
	}

	var result []map[string]string
	if err = json.Unmarshal(resp, &result); err != nil {
		return "", err
	}

	return result[0]["translation_text"], nil
}
