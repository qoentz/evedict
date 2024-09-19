package server

import (
	"errors"
	"evedict/config"
	"evedict/internal/handler"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
)

func initialize(config *config.SystemConfig) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	api := router.PathPrefix("/api").Subrouter()

	api.Handle("/news", handler.GetNews(config.PromptTemplate)).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusNotFound)
	})

	return router

}

func ServeRouter(config *config.SystemConfig) *http.Server {
	r := initialize(config)

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
