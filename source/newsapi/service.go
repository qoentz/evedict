package newsapi

import (
	"fmt"
	"os"
	"strings"
)

func CreatePromptFromHeadlines(articles []Article) string {
	var headlines []string
	for i, article := range articles {
		if i >= 5 {
			break
		}

		if article.Content != "" {
			headlines = append(headlines, fmt.Sprintf("Title: %s\nContent: %s", article.Title, article.Content))
		}
	}

	prompt := os.Getenv("PROMPT_INPUT")
	prompt += strings.Join(headlines, "\n- ")

	return prompt
}
