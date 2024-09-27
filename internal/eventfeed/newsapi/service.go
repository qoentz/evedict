package newsapi

import (
	"encoding/json"
	"github.com/qoentz/evedict/internal/httputil"
	"io"
)

type Service struct {
	APIKey  string
	BaseURL string
}

func NewNewsAPIService(apiKey, baseURL string) *Service {
	return &Service{
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
}

func (n *Service) Fetch() ([]Article, error) {
	resp, err := httputil.GetRequest(n.BaseURL + n.APIKey)
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
