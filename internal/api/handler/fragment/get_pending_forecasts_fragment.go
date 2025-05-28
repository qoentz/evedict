package fragment

import (
	"fmt"
	"github.com/qoentz/evedict/internal/service"
	"github.com/qoentz/evedict/internal/util"
	"github.com/qoentz/evedict/internal/view"
	"net/http"
	"strconv"
)

func GetPendingForecastsFragment(s *service.ForecastService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		limit := 2
		offset := 0
		if o := r.URL.Query().Get("offset"); o != "" {
			var err error
			offset, err = strconv.Atoi(o)
			if err != nil {
				http.Error(w, "invalid offset", http.StatusBadRequest)
				return
			}
		}

		var category *util.Category
		if catStr := r.URL.Query().Get("category"); catStr != "" {
			cat, err := util.ParseCategory(catStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			category = &cat
		}

		// Get forecasts + hasMore
		forecasts, hasMore, err := s.GetPendingForecasts(limit, offset, category)
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't get forecasts: %v", err), http.StatusInternalServerError)
			return
		}

		// Render the section
		err = view.PendingForecastSection(forecasts, offset, hasMore).Render(r.Context(), w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
