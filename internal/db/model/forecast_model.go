package model

import (
	"github.com/google/uuid"
	"github.com/qoentz/evedict/internal/util"
	"time"
)

type Forecast struct {
	ID        uuid.UUID     `db:"id"`
	Headline  string        `db:"headline"`
	Summary   string        `db:"summary"`
	ImageURL  string        `db:"image_url"`
	Category  util.Category `db:"category"`
	Timestamp time.Time     `db:"timestamp"`
	Outcomes  []Outcome     `db:"-"`
	Tags      []Tag         `db:"-"`
	Sources   []Source      `db:"-"`
	Market    *Market       `db:"-"`
}
