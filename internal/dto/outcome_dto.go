package dto

import "github.com/google/uuid"

type Outcome struct {
	ID              uuid.UUID `json:"id"`
	Content         string    `json:"content"`
	ConfidenceLevel int       `json:"confidenceLevel"`
}
