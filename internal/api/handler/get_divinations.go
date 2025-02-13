package handler

import (
	"fmt"
	"github.com/qoentz/evedict/internal/service"
	"github.com/qoentz/evedict/internal/view"
	"net/http"
)

func GetDivinations(s *service.DivinationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		divinations, err := s.GetDivinations()
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't get divinations: %v", err), http.StatusInternalServerError)
			return
		}

		err = view.DivinationFeed(divinations).Render(r.Context(), w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
