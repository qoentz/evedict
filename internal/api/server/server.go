package server

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/qoentz/evedict/internal/api/handler"
	"github.com/qoentz/evedict/internal/api/handler/fragment"
	"github.com/qoentz/evedict/internal/api/handler/page"
	"github.com/qoentz/evedict/internal/api/middleware"
	"github.com/qoentz/evedict/internal/registry"
	"github.com/qoentz/evedict/internal/view"
	"log"
	"net"
	"net/http"
)

func InitRouter(reg *registry.Registry) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Public routes — no auth
	router.HandleFunc("/login", page.Login()).Methods("GET")
	router.HandleFunc("/login", handler.SubmitPassword(reg.AuthService)).Methods("POST")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/favicon.ico")
	}).Methods("GET")

	// Full page endpoints
	router.HandleFunc("/", page.Home()).Methods("GET")
	router.HandleFunc("/forecasts/{forecastId}", page.GetForecast(reg.ForecastService)).Methods("GET")
	router.HandleFunc("/about", page.About()).Methods("GET")
	router.HandleFunc("/contact", page.Contact()).Methods("GET")

	// API endpoints for htmx partial updates
	api := router.PathPrefix("/api").Subrouter()
	api.Handle("/forecasts", fragment.GetForecastsFragment(reg.ForecastService)).Methods("GET")
	api.Handle("/forecasts/{forecastId}", fragment.GetForecastFragment(reg.ForecastService)).Methods("GET")
	api.Handle("/about", fragment.AboutFragment()).Methods("GET")
	api.Handle("/contact", fragment.ContactFragment()).Methods("GET")

	// Protected subrouter
	protected := router.NewRoute().Subrouter()
	protected.Use(middleware.AuthShield(reg.AuthService)) // <- wrap everything below

	// Admin endpoints
	vault := protected.PathPrefix("/vault").Subrouter()
	vault.HandleFunc("/workspace", page.WorkSpace()).Methods("GET")
	vault.HandleFunc("/workspace/pending", fragment.GetPendingForecastsFragment(reg.ForecastService)).Methods("GET")
	vault.HandleFunc("/forecasts/{forecastId}", handler.ApproveForecast(reg.ForecastService)).Methods("PATCH")

	invoke := vault.PathPrefix("/invoke").Subrouter()
	invoke.Handle("/forecast/default", handler.GenerateForecasts(reg.ForecastService)).Methods("POST")
	invoke.Handle("/forecast/poly", handler.GeneratePolyForecasts(reg.ForecastService)).Methods("POST")

	// Not found
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err := view.NotFound().Render(r.Context(), w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	})

	return router
}

func ServeHTTP(r *mux.Router) *http.Server {
	serverAddr := fmt.Sprintf("0.0.0.0:%s", "8080")

	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", serverAddr, err)
	}

	log.Printf("Serving HTTP at %s", serverAddr)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server stopped listening on %s: %v", serverAddr, err)
		}
	}()

	return server
}
