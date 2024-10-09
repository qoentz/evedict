package model

import "github.com/google/uuid"

type Outcome struct {
	ID              uuid.UUID `db:"id"`
	PredictionID    uuid.UUID `db:"prediction_id"`
	Content         string    `db:"content"`
	ConfidenceLevel int       `db:"confidence_level"`
}
