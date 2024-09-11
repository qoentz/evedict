package newsapi

import (
	"os"
)

func CreatePromptFromHeadlines(articles []Article) string {
	var headlines []string
	for _, article := range articles {
		headlines = append(headlines, article.Content)
	}

	prompt := os.Getenv("PROMPT_INPUT")
	//prompt += strings.Join(headlines, "\n- ")

	//fmt.Println(prompt)

	return prompt
}
