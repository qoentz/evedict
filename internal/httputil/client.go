package httputil

import (
	"bytes"
	"io"
	"net/http"
)

func PostRequest(url string, reqBody []byte, apiKey string) (io.ReadCloser, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	//if resp.StatusCode != http.StatusOK {
	//	return resp.Body, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	//}

	return resp.Body, nil
}

func GetRequest(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	//if resp.StatusCode != http.StatusOK {
	//	return resp.Body, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	//}

	return resp.Body, nil
}
