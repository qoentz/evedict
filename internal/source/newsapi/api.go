package newsapi

import (
	"encoding/json"
	"evedict/internal/httputil"
	"io"
	"os"
)

func Fetch() ([]Article, error) {
	resp, err := httputil.GetRequest(os.Getenv("NEWS_API_URL") + os.Getenv("NEWS_API_KEY"))
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	respBody, err := io.ReadAll(resp)
	if err != nil {
		return nil, err
	}

	var data Response
	if err = json.Unmarshal(respBody, &data); err != nil {
		return nil, err
	}

	return data.Articles, nil
}
