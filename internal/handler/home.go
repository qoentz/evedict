package handler

import (
	"fmt"
	"github.com/qoentz/evedict/internal/view"
	"net/http"
)

func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err := view.Home().Render(r.Context(), w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering main site: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
