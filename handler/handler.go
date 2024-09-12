package handler

import (
	"encoding/json"
	"evedict/llm/replicate"
	"evedict/source/newsapi"
	"fmt"
	"net/http"
)

func GetNews(w http.ResponseWriter, r *http.Request) {
	data, err := newsapi.Fetch()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data from GDELT: %v", err), http.StatusInternalServerError)
		return
	}

	url, err := replicate.InitiateStream(newsapi.CreatePromptFromHeadlines(data))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error initiating stream: %v", err), http.StatusInternalServerError)
		return
	}

	predictions, err := replicate.HandleStream(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing stream: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(predictions)
	if err != nil {
		return
	}
}
