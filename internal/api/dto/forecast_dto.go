package dto

import (
	"github.com/google/uuid"
	"time"
)

type Forecast struct {
	ID        uuid.UUID `json:"id"`
	Headline  string    `json:"headline"`
	Summary   string    `json:"summary"`
	Outcomes  []Outcome `json:"outcomes"`
	ImageURL  string    `json:"imageUrl"`
	Sources   []Source  `json:"sources"`
	Timestamp time.Time `json:"timestamp"`
}
