package main

import (
	"evedict/llm/replicate"
	"evedict/source/newsapi"
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	data, err := newsapi.Fetch()
	if err != nil {
		log.Fatalf("Error fetching data from GDELT: %v", err)
	}

	url, err := replicate.InitiateStream(newsapi.CreatePromptFromHeadlines(data))
	if err != nil {
		log.Fatalf("Error initiating stream: %v", err)
	}

	err = replicate.HandleStream(url)
	if err != nil {
		fmt.Println("Error handling stream:", err)
		return
	}
}
