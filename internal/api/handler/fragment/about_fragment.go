package fragment

import (
	"fmt"
	"github.com/qoentz/evedict/internal/view"
	"net/http"
)

func AboutFragment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err := view.AboutFragment().Render(r.Context(), w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
