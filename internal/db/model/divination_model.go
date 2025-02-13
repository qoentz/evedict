package model

import (
	"github.com/google/uuid"
	"time"
)

type Divination struct {
	ID        uuid.UUID `db:"id"`
	Headline  string    `db:"headline"`
	Summary   string    `db:"summary"`
	ImageURL  string    `db:"image_url"`
	Timestamp time.Time `db:"timestamp"`
	Outcomes  []Outcome `db:"-"`
	Sources   []Source  `db:"-"`
}
