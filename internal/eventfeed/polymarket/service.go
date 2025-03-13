package polymarket

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Service struct {
	HTTPClient *http.Client
	BaseURL    string
}

func NewPolyMarketService(client *http.Client, baseURL string) *Service {
	return &Service{
		HTTPClient: client,
		BaseURL:    baseURL,
	}
}

func (s *Service) FetchTopEvents() ([]Event, error) {
	weekAgo := time.Now().UTC().AddDate(0, 0, -7)
	startDate := weekAgo.Format("2006-01-02T15:04:05Z")

	url := fmt.Sprintf(
		"%s/events?start_date_min=%s&volume_min=5000&closed=false",
		s.BaseURL,
		startDate,
	)

	events, err := s.Fetch(url)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Service) Fetch(url string) ([]Event, error) {
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

	var data []Event
	if err = json.Unmarshal(respBody, &data); err != nil {
		return nil, err
	}

	return data, nil
}
