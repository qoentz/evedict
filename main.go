package main

import (
	"context"
	"evedict/server"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	httpServer := server.ServeRouter()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	//data, err := newsapi.Fetch()
	//if err != nil {
	//	log.Fatalf("Error fetching data from GDELT: %v", err)
	//}
	//
	//url, err := replicate.InitiateStream(newsapi.CreatePromptFromHeadlines(data))
	//if err != nil {
	//	log.Fatalf("Error initiating stream: %v", err)
	//}
	//
	//err = replicate.HandleStream(url)
	//if err != nil {
	//	fmt.Println("Error handling stream:", err)
	//	return
	//}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
}
