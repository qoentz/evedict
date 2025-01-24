package handler

import (
	"fmt"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"github.com/qoentz/evedict/internal/service"
	"github.com/qoentz/evedict/internal/view"
	"net/http"
)

func GeneratePredictions(s *service.PredictionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		category := r.URL.Query().Get("category")
		if category == "" {
			category = "general"
		}

		newsCategory, err := newsapi.ValidateCategory(category)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid category: %v", err), http.StatusBadRequest)
			return
		}

		predictions, err := s.GeneratePredictions(newsCategory)
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't generate predictions: %v", err), http.StatusInternalServerError)
			return
		}

		err = s.SavePredictions(predictions)
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't save predictions: %v", err), http.StatusInternalServerError)
			return
		}

		err = view.HighlightedSlider(predictions).Render(r.Context(), w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}

		//w.Header().Set("Content-Type", "application/json")
		//err = json.NewEncoder(w).Encode(predictions)
		//if err != nil {
		//	return
		//}
	}
}
