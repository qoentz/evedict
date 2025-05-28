package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/qoentz/evedict/internal/service"
	"net/http"
)

func ApproveForecast(s *service.ForecastService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		forecastID, err := uuid.Parse(mux.Vars(r)["forecastId"])
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid forecast ID: %v", err), http.StatusBadRequest)
			return
		}

		if err = s.ApproveForecast(forecastID); err != nil {
			http.Error(w, "Could not approve forecast: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
