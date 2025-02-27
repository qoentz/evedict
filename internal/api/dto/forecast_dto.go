package dto

import (
	"github.com/google/uuid"
	"github.com/qoentz/evedict/internal/util"
	"time"
)

type Forecast struct {
	ID        uuid.UUID     `json:"id"`
	Headline  string        `json:"headline"`
	Summary   string        `json:"summary"`
	Outcomes  []Outcome     `json:"outcomes"`
	Category  util.Category `json:"category"`
	ImageURL  string        `json:"imageUrl"`
	Tags      []Tag         `json:"tags"`
	Sources   []Source      `json:"sources"`
	Timestamp time.Time     `json:"timestamp"`
	Related   []Forecast    `json:"related"`
}
