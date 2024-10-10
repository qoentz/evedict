package gdelt

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func Fetch() ([]Article, error) {
	resp, err := http.Get(os.Getenv("GDELT_URL"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data Response
	if err = json.Unmarshal(resBody, &data); err != nil {
		return nil, err
	}

	return data.Articles, nil
}
