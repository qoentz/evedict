package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/qoentz/evedict/internal/service"
	"github.com/qoentz/evedict/internal/view"
	"net/http"
)

func GetForecastFragment(s *service.ForecastService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		forecastId, err := uuid.Parse(mux.Vars(r)["forecastId"])
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid forecast ID: %v", err), http.StatusBadRequest)
			return
		}

		forecast, err := s.GetForecast(forecastId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't get forecast: %v", err), http.StatusInternalServerError)
			return
		}

		forecasts, err := s.GetForecasts(4, 9)
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't get forecasts: %v", err), http.StatusInternalServerError)
			return
		}

		err = view.ForecastDetailFragment(forecast, forecasts).Render(r.Context(), w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
