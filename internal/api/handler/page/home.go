package page

import (
	"fmt"
	"github.com/qoentz/evedict/internal/util"
	"github.com/qoentz/evedict/internal/view"
	"net/http"
)

func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		query := r.URL.Query()

		var category *util.Category
		if catStr := query.Get("category"); catStr != "" {
			cat, err := util.ParseCategory(catStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			category = &cat
		}

		err := view.Home(category).Render(r.Context(), w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering main site: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
