package fragment

import (
	"fmt"
	"github.com/qoentz/evedict/internal/service"
	"github.com/qoentz/evedict/internal/util"
	"github.com/qoentz/evedict/internal/view"
	"net/http"
)

func GetForecastsFragment(s *service.ForecastService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		query := r.URL.Query()
		limit, offset, err := util.ParsePagination(query, 9, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var category *util.Category
		if catStr := query.Get("category"); catStr != "" {
			cat, err := util.ParseCategory(catStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			category = &cat
		}

		forecasts, err := s.GetForecasts(limit, offset, category)
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
