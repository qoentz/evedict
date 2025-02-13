package handler

import (
	"fmt"
	"github.com/qoentz/evedict/internal/service"
	"github.com/qoentz/evedict/internal/view"
	"net/http"
)

func GetForecasts(s *service.ForecastService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		forecasts, err := s.GetForecasts()
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't get forecasts: %v", err), http.StatusInternalServerError)
			return
		}

		err = view.ForecastFeed(forecasts).Render(r.Context(), w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
