package handler

import (
	"encoding/json"
	"fmt"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"github.com/qoentz/evedict/internal/llm"
	"github.com/qoentz/evedict/internal/promptgen"
	"net/http"
)

func GetNews(newsAPI *newsapi.Service, ai llm.Service, template *promptgen.PromptTemplate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := newsAPI.Fetch()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching data from GDELT: %v", err), http.StatusInternalServerError)
			return
		}

		prompt, err := template.CreatePromptWithArticles(data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error building prompt: %v", err), http.StatusInternalServerError)
			return
		}

		predictions, err := ai.GetPredictions(prompt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting response: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(predictions)
		if err != nil {
			return
		}
	}
}
