package handler

import (
	"fmt"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"github.com/qoentz/evedict/internal/llm"
	"github.com/qoentz/evedict/internal/promptgen"
	"net/http"
)

func GetBusinessPredictions(newsAPI *newsapi.Service, ai llm.Service, template *promptgen.PromptTemplate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := newsAPI.FetchTopHeadlines(newsapi.Business)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching data from NewsAPI: %v", err), http.StatusInternalServerError)
			return
		}

		prompt, err := template.CreateKeywordExtractionPrompt(data[0])
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating keyword extraction prompt: %v", err), http.StatusInternalServerError)
			return
		}

		keywords, err := ai.ExtractKeywords(prompt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error extracting keywords: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Println(keywords)

		//w.Header().Set("Content-Type", "application/json")
		//err = json.NewEncoder(w).Encode(data)
		//if err != nil {
		//	return
		//}
	}
}
