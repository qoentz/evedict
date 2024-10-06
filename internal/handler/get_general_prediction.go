package handler

import (
	"fmt"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"github.com/qoentz/evedict/internal/llm"
	"github.com/qoentz/evedict/internal/promptgen"
	"github.com/qoentz/evedict/internal/view"
	"net/http"
)

func GetGeneralPredictions(newsAPI *newsapi.Service, ai llm.Service, template *promptgen.PromptTemplate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headlines, err := newsAPI.FetchTopHeadlines(newsapi.General)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching headlines from NewsAPI: %v", err), http.StatusInternalServerError)
			return
		}

		selectionPrompt, err := template.CreateArticleSelectionPrompt(headlines)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating article selection prompt: %v", err), http.StatusInternalServerError)
			return
		}

		articleSelection, err := ai.SelectArticles(selectionPrompt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error selecting articles: %v", err), http.StatusInternalServerError)
			return
		}

		var predictions []llm.Prediction
		for _, idx := range articleSelection {
			mainArticle := headlines[idx]

			extractionPrompt, err := template.CreateKeywordExtractionPrompt(mainArticle)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error creating keyword extraction prompt: %v", err), http.StatusInternalServerError)
				return
			}

			keywords, err := ai.ExtractKeywords(extractionPrompt)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error extracting keywords: %v", err), http.StatusInternalServerError)
				return
			}

			articles, err := newsAPI.FetchWithKeywords(keywords)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error fetching articles from NewsAPI with keywords: %v", err), http.StatusInternalServerError)
				return
			}

			predictionPrompt, err := template.CreatePredictionPrompt(mainArticle, articles)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error creating prediction prompt: %v", err), http.StatusInternalServerError)
				return
			}

			prediction, err := ai.GetPrediction(predictionPrompt)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error generating prediction: %v", err), http.StatusInternalServerError)
				return
			}

			prediction.ImageURL = mainArticle.URLToImage
			predictions = append(predictions, *prediction)
		}

		//predictions := []llm.Prediction{
		//	{
		//		Headline: "Tropical Storm Milton to Impact Florida",
		//		Summary:  "Tropical Storm Milton is forecast to strengthen into a hurricane and head towards Florida, potentially impacting the western coast.",
		//		Outcomes: []llm.Outcome{
		//			{
		//				Content:         "Milton makes landfall in Florida as a hurricane, causing significant damage and disruption.",
		//				ConfidenceLevel: 85,
		//			},
		//			{
		//				Content:         "Milton weakens before reaching Florida, resulting in minimal damage and impact",
		//				ConfidenceLevel: 40,
		//			},
		//		},
		//		ImageURL: "https://assets1.cbsnewsstatic.com/hub/i/r/2024/10/05/d93be9e3-3d8a-465b-b7c3-12e046e601d6/thumbnail/1200x630/2cc76f942edcd5b45a3ff75d6ae1fb8b/milton.jpg?v=0736ad3ef1e9ddfe1218648fe91d6c9b",
		//	},
		//}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = view.PredictionFeed(predictions).Render(r.Context(), w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
