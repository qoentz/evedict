package model

import (
	"github.com/google/uuid"
	"time"
)

type RelatedForecast struct {
	ID           uuid.UUID `db:"id"`
	Headline     string    `db:"headline"`
	Summary      string    `db:"summary"`
	ImageURL     string    `db:"image_url"`
	Timestamp    time.Time `db:"timestamp"`
	MatchedByTag bool      `db:"matched_by_tag"`
}
