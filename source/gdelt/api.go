package gdelt

import (
	"cashnew/httputil"
	"encoding/json"
	"os"
)

func Fetch() ([]Article, error) {
	resp, err := httputil.GetRequest(os.Getenv("GDELT_URL"))
	if err != nil {
		return nil, err
	}

	var data Response
	if err = json.Unmarshal(resp, &data); err != nil {
		return nil, err
	}

	return data.Articles, nil
}
