package model

import "github.com/google/uuid"

type Outcome struct {
	ID              uuid.UUID `db:"id"`
	ForecastID      uuid.UUID `db:"forecast_id"`
	Content         string    `db:"content"`
	ConfidenceLevel int       `db:"confidence_level"`
}
