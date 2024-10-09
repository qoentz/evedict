package handler

import (
	"fmt"
	"github.com/qoentz/evedict/internal/service"
	"github.com/qoentz/evedict/internal/view"
	"net/http"
)

func GetPredictions(s *service.PredictionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		predictions, err := s.GetPredictions()
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't get predictions: %v", err), http.StatusInternalServerError)
			return
		}

		err = view.PredictionFeed(predictions).Render(r.Context(), w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
