package handler

import (
	"fmt"
	"github.com/qoentz/evedict/internal/view"
	"net/http"
)

func MainFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set content type to HTML
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Render the main site with the templ component
		err := view.MainFeed().Render(r.Context(), w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering main site: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
