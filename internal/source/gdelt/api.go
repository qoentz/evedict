package gdelt

import (
	"encoding/json"
	"evedict/internal/httputil"
	"io"
	"os"
)

func Fetch() ([]Article, error) {
	resp, err := httputil.GetRequest(os.Getenv("GDELT_URL"))
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	resBody, err := io.ReadAll(resp)
	if err != nil {
		return nil, err
	}

	var data Response
	if err = json.Unmarshal(resBody, &data); err != nil {
		return nil, err
	}

	return data.Articles, nil
}
