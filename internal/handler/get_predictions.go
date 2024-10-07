package handler

import (
	"fmt"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"github.com/qoentz/evedict/internal/llm"
	"github.com/qoentz/evedict/internal/promptgen"
	"github.com/qoentz/evedict/internal/view"
	"net/http"
)

func GetPredictions(newsAPI *newsapi.Service, ai llm.Service, template *promptgen.PromptTemplate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := r.URL.Query().Get("category")
		if category == "" {
			category = "general"
		}

		newsCategory, err := newsAPI.ParseCategory(category)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid category: %v", err), http.StatusBadRequest)
			return
		}

		headlines, err := newsAPI.FetchTopHeadlines(newsCategory)
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

			var sources []llm.Source

			mainSource := llm.Source{
				Name:  mainArticle.Source.Name,
				Title: mainArticle.Title,
				URL:   mainArticle.URL,
			}
			sources = append(sources, mainSource)

			for _, article := range articles {
				if article.URL == mainArticle.URL {
					continue
				}

				source := llm.Source{
					Name:  article.Source.Name,
					Title: article.Title,
					URL:   article.URL,
				}

				sources = append(sources, source)
			}

			prediction.Sources = sources
			predictions = append(predictions, *prediction)
		}

		//predictions := []llm.Prediction{
		//	{
		//		Headline: "CM Punk's Victory at WWE Bad Blood 2024 Sets Stage for Future Rivalries",
		//		Summary:  "With CM Punk's win over Drew McIntyre at WWE Bad Blood 2024, the stage is set for future grudge matches and potential rivalries, including a possible feud with Judgment Day.",
		//		Outcomes: []llm.Outcome{
		//			{
		//				Content:         "CM Punk will face Damian Priest in a future premium live event, reigniting their rivalry and leading to a potential title shot.",
		//				ConfidenceLevel: 85,
		//			},
		//			{
		//				Content:         "Nia Jax's retention of the WWE Women's Championship will lead to a prolonged feud with Bayley, with Tiffany Stratton potentially playing a key role in the storyline.",
		//				ConfidenceLevel: 70,
		//			},
		//		},
		//		ImageURL: "https://sportshub.cbsistatic.com/i/r/2024/10/03/15b4e4b1-860c-4c7b-a10f-52ff3c4a3544/thumbnail/1200x675/ccbcc8e93f34a6b94258f1b400b7c4f9/cm-punk-cage.jpg",
		//		Sources: []llm.Source{
		//			{
		//				Title: "M4 MacBook Pro: Four things to expect with Apple’s next Pro laptop",
		//				Name:  "9to5Mac",
		//				URL:   "https://9to5mac.com/2024/10/06/m4-macbook-pro-roundup/",
		//			},
		//		},
		//	},
		//	{
		//		Headline: "Joker 2's Box Office Performance to Suffer Due to Poor Reception",
		//		Summary:  "The sequel to the critically acclaimed Joker film has received a historically low CinemaScore, indicating a poor reception from audiences, which may negatively impact its box office performance.",
		//		Outcomes: []llm.Outcome{
		//			{
		//				Content:         "Joker 2 will fail to reach the $200 million mark in its domestic box office run, leading to a significant financial loss for the studio.",
		//				ConfidenceLevel: 80,
		//			},
		//			{
		//				Content:         "Despite the poor reception, Joker 2 will still manage to gross over $300 million domestically, but its performance will be considered a disappointment compared to its predecessor.",
		//				ConfidenceLevel: 60,
		//			},
		//		},
		//		ImageURL: "https://variety.com/wp-content/uploads/2024/04/Screen-Shot-2024-04-09-at-9.01.48-PM.png?w=1000&h=563&crop=1",
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
