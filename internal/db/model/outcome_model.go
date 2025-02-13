package model

import "github.com/google/uuid"

type Outcome struct {
	ID              uuid.UUID `db:"id"`
	DivinationID    uuid.UUID `db:"divination_id"`
	Content         string    `db:"content"`
	ConfidenceLevel int       `db:"confidence_level"`
}
