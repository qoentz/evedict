package util

import (
	"fmt"
	"github.com/qoentz/evedict/internal/eventfeed/newsapi"
	"strings"
)

type Category string

const (
	Politics   Category = "Politics"
	Economy    Category = "Economy"
	Technology Category = "Technology"
	Culture    Category = "Culture"
)

var validCategories = map[string]Category{
	"Politics":   Politics,
	"Economy":    Economy,
	"Technology": Technology,
	"Culture":    Culture,
}

func ParseCategory(s string) (Category, error) {
	for key, cat := range validCategories {
		if strings.EqualFold(s, key) {
			return cat, nil
		}
	}
	return "", fmt.Errorf("invalid category: %s", s)
}

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
