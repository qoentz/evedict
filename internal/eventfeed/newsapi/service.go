package newsapi

import (
	"encoding/json"
	"io"
	"net/http"
)

type Service struct {
	HTTPClient *http.Client
	APIKey     string
	BaseURL    string
}

func NewNewsAPIService(client *http.Client, apiKey, baseURL string) *Service {
	return &Service{
		HTTPClient: client,
		APIKey:     apiKey,
		BaseURL:    baseURL,
	}
}

func (n *Service) Fetch() ([]Article, error) {
	req, err := http.NewRequest("GET", n.BaseURL+n.APIKey, nil)
	if err != nil {
		return nil, err
	}

	resp, err := n.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data Response
	if err = json.Unmarshal(respBody, &data); err != nil {
		return nil, err
	}

	return data.Articles, nil
}
