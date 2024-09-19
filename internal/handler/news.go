package handler

import (
	"encoding/json"
	"evedict/internal/llm/replicate"
	"evedict/internal/promptgen"
	"evedict/internal/source/newsapi"
	"fmt"
	"net/http"
)

func GetNews(promptTemplate *promptgen.PromptTemplate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := newsapi.Fetch()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching data from GDELT: %v", err), http.StatusInternalServerError)
			return
		}

		prompt, err := promptTemplate.CreatePromptWithArticles(data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error building prompt: %v", err), http.StatusInternalServerError)
			return
		}

		predictions, err := replicate.GetPredictions(prompt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting response: %v", err), http.StatusInternalServerError)
			return
		}

		//url, err := replicate.InitiateStream(prompt)
		//if err != nil {
		//	http.Error(w, fmt.Sprintf("Error initiating stream: %v", err), http.StatusInternalServerError)
		//	return
		//}
		//
		//predictions, err := replicate.HandleStream(url)
		//if err != nil {
		//	http.Error(w, fmt.Sprintf("Error processing stream: %v", err), http.StatusInternalServerError)
		//	return
		//}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(predictions)
		if err != nil {
			return
		}
	}
}
