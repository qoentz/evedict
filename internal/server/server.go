package server

import (
	"errors"
	"evedict/config"
	"evedict/internal/handler"
	"evedict/internal/registry"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
)

func InitRouter(config *config.SystemConfig, reg *registry.Registry) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	api := router.PathPrefix("/api").Subrouter()

	api.Handle("/news", handler.GetNews(reg.NewsAPIService, reg.ReplicateService, config.PromptTemplate)).Methods("GET")

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
