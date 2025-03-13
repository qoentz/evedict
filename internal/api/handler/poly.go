package handler

import (
	"encoding/json"
	"fmt"
	"github.com/qoentz/evedict/internal/service"
	"net/http"
)

func GeneratePolyForecasts(s *service.ForecastService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		forecasts, err := s.GeneratePolyForecasts()
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't generate forecasts: %v", err), http.StatusInternalServerError)
			return
		}

		err = s.SavePolyForecasts(forecasts)
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
