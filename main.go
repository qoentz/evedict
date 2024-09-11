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

	sum, err := replicate.GenerateSummary(newsapi.CreatePromptFromHeadlines(data))
	if err != nil {
		log.Fatalf("Error generating summary: %v", err)
	}

	fmt.Println(sum)

	//getURL := "https://api.replicate.com/v1/predictions/ddn8k254ssrj20chvttr8jv264"
	//
	//// Check the status of the prediction
	//predictionResult, err := replicate.CheckPredictionStatus(getURL)
	//if err != nil {
	//	log.Fatalf("Error getting prediction result: %v", err)
	//}
	//
	//// Format the prediction nicely
	//formattedPrediction := replicate.FormatPrediction(predictionResult)
	//
	//// Print the formatted prediction result
	//fmt.Println(formattedPrediction)
}
