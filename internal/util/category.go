package util

import "github.com/qoentz/evedict/internal/eventfeed/newsapi"

type Category string

const (
	Politics   Category = "Politics"
	Economy    Category = "Economy"
	Technology Category = "Technology"
	Culture    Category = "Culture"
)

// DetermineCategory Might be more useful than LLM
func DetermineCategory(category newsapi.Category) Category {
	switch category {
	case newsapi.Business:
		return Economy
	case newsapi.Entertainment, newsapi.Sports:
	default:
		return Politics
	}

	return ""
}
