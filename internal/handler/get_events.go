package handler

import (
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"github.com/qoentz/evedict/internal/llm"
	"github.com/qoentz/evedict/internal/promptgen"
	"net/http"
)

func GetEvents(newsAPI *newsapi.Service, ai llm.Service, template *promptgen.PromptTemplate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//data, err := newsAPI.Fetch("e")
		//if err != nil {
		//	http.Error(w, fmt.Sprintf("Error fetching data from GDELT: %v", err), http.StatusInternalServerError)
		//	return
		//}
		//
		//prompt, err := template.CreatePredictionPrompt(data)
		//if err != nil {
		//	http.Error(w, fmt.Sprintf("Error building prompt: %v", err), http.StatusInternalServerError)
		//	return
		//}
		//
		//predictions, err := ai.GetGeneralPredictions(prompt, data)
		//if err != nil {
		//	http.Error(w, fmt.Sprintf("Error getting response: %v", err), http.StatusInternalServerError)
		//	return
		//}

		//predictions := &llm.Predictions{
		//	Predictions: []llm.Prediction{
		//		{
		//			Title:   "Boar's Head Listeria Outbreak Investigation Concludes",
		//			Content: "The investigation into the Boar's Head listeria outbreak that killed 10 people will conclude and results will be made public.",
		//		},
		//		{
		//			Title:   "San Jose, California Approves Sixth Costco Location",
		//			Content: "The San Jose city council will approve the construction of a sixth Costco location in the city, making it the first US city to have six Costco stores.",
		//		},
		//	},
		//}

		//w.Header().Set("Content-Type", "text/html; charset=utf-8")
		//err = view.EventFeed(predictions.Predictions).Render(r.Context(), w)
		//if err != nil {
		//	http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		//	return
		//}

		//w.Header().Set("Content-Type", "application/json")
		//err = json.NewEncoder(w).Encode(predictions)
		//if err != nil {
		//	return
		//}
	}
}
