package newsapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

func (s *Service) FetchTopHeadlines(category Category) ([]Article, error) {
	params := map[string]string{
		"category": string(category),
	}

	path, err := s.ConstructURL(TopHeadlines, params)
	if err != nil {
		return nil, err
	}

	articles, err := s.Fetch(path)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (s *Service) FetchWithKeywords(keywords []string) ([]Article, error) {
	if len(keywords) == 0 {
		return nil, fmt.Errorf("no keywords provided")
	}

	query := strings.Join(keywords, " ")

	params := map[string]string{
		"q":        query,
		"pageSize": "10",
		"sortBy":   "publishedAt",
	}

	path, err := s.ConstructURL(Everything, params)
	if err != nil {
		return nil, err
	}

	articles, err := s.Fetch(path)
	if err != nil {
		return nil, err
	}

	return articles, nil

}

func (s *Service) ConstructURL(endpoint Endpoint, params map[string]string) (string, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%s", s.BaseURL, endpoint))
	if err != nil {
		return "", err
	}

	query := u.Query()

	for key, value := range params {
		query.Set(key, value)
	}

	query.Set("apiKey", s.APIKey)

	u.RawQuery = query.Encode()

	return u.String(), nil
}

func (s *Service) Fetch(url string) ([]Article, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("API returned status: %s, but the error message could not be read", resp.Status)
		}

		return nil, fmt.Errorf("API returned status: %s, message: %s", resp.Status, string(respBody))
	}

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
