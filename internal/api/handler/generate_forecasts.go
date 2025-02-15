package handler

import (
	"encoding/json"
	"fmt"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"github.com/qoentz/evedict/internal/service"
	"net/http"
)

func GenerateForecasts(s *service.ForecastService) http.HandlerFunc {
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

		forecasts, err := s.GenerateForecasts(newsCategory)
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't generate forecasts: %v", err), http.StatusInternalServerError)
			return
		}

		err = s.SaveForecasts(forecasts)
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't save forecasts: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(forecasts)
		if err != nil {
			return
		}
	}
}
