package handler

import (
	"fmt"
	"github.com/qoentz/evedict/internal/service"
	"github.com/qoentz/evedict/internal/view"
	"net/http"
)

func SubmitPassword(auth *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form", http.StatusBadRequest)
			return
		}

		if r.FormValue("password") == auth.AuthSecret {
			auth.IssueToken(w)
			w.Header().Set("HX-Redirect", "/")
			return
		}

		fmt.Println("WRONG")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = view.LoginForm(true).Render(r.Context(), w)
	}
}
