package server

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/qoentz/evedict/config"
	"github.com/qoentz/evedict/internal/handler"
	"github.com/qoentz/evedict/internal/registry"
	"log"
	"net"
	"net/http"
)

func InitRouter(config *config.SystemConfig, reg *registry.Registry) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", handler.Home()).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.FileServer(http.Dir("./static/")).ServeHTTP(w, r)
	})))

	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/favicon.ico")
	}).Methods("GET")

	api := router.PathPrefix("/api").Subrouter()
	api.Handle("/predictions", handler.GetPredictions(reg.NewsAPIService, reg.ReplicateService, config.PromptTemplate)).Methods("GET")

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
