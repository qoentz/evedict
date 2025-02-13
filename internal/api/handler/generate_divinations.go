package handler

import (
	"fmt"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"github.com/qoentz/evedict/internal/service"
	"github.com/qoentz/evedict/internal/view"
	"net/http"
)

func GenerateDivinations(s *service.DivinationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		category := r.URL.Query().Get("category")
		if category == "" {
			category = "general"
		}

		newsCategory, err := newsapi.ValidateCategory(category)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid category: %v", err), http.StatusBadRequest)
			return
		}

		divinations, err := s.GenerateDivinations(newsCategory)
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't generate divinations: %v", err), http.StatusInternalServerError)
			return
		}

		err = s.SaveDivinations(divinations)
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't save divinations: %v", err), http.StatusInternalServerError)
			return
		}

		err = view.DivinationFeed(divinations).Render(r.Context(), w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}

		//w.Header().Set("Content-Type", "application/json")
		//err = json.NewEncoder(w).Encode(divinations)
		//if err != nil {
		//	return
		//}
	}
}
