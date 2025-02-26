package server

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/qoentz/evedict/internal/api/handler"
	"github.com/qoentz/evedict/internal/registry"
	"log"
	"net"
	"net/http"
)

func InitRouter(reg *registry.Registry) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Full page endpoints
	router.HandleFunc("/", handler.Home()).Methods("GET")
	router.HandleFunc("/forecasts/{forecastId}", handler.GetForecast(reg.ForecastService)).Methods("GET")

	// Static assets and favicon
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/favicon.ico")
	}).Methods("GET")

	// API endpoints for htmx partial updates
	api := router.PathPrefix("/api").Subrouter()
	api.Handle("/forecasts", handler.GetForecasts(reg.ForecastService)).Methods("GET")
	api.Handle("/forecasts/{forecastId}", handler.GetForecastFragment(reg.ForecastService)).Methods("GET")

	api.Handle("/gen", handler.GenerateForecasts(reg.ForecastService)).Methods("POST")

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusNotFound)
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
