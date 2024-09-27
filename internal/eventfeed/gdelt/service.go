package gdelt

import (
	"fmt"
	"github.com/qoentz/evedict/internal/llm/hugface"
	"log"
	"os"
	"strings"
)

func CreatePromptFromHeadlines(articles []Article) string {
	var headlines []string
	for _, article := range articles {
		if article.Language != "English" {
			translatedTitle, err := hugface.TranslateHeadline(article.Title)
			if err != nil {
				log.Printf("Error translating headline: %v", err)
				continue
			}
			fmt.Println(translatedTitle)
			headlines = append(headlines, translatedTitle)
		} else {
			fmt.Println(article.Title)
			headlines = append(headlines, article.Title)
		}
	}

	prompt := os.Getenv("PROMPT_INPUT")
	prompt += strings.Join(headlines, "\n- ")

	return prompt
}
