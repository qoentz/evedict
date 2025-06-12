package middleware

import (
	"github.com/gorilla/mux"
	"github.com/qoentz/evedict/internal/service"
	"net/http"
)

func AuthShield(auth *service.AuthService) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth")
			if err != nil || !auth.ValidateToken(cookie.Value) {
				if r.Header.Get("HX-Request") != "" {
					w.Header().Set("HX-Redirect", "/vault/login")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				http.Redirect(w, r, "/vault/login", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
