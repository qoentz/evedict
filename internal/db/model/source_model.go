package model

import "github.com/google/uuid"

type Source struct {
	ID           uuid.UUID `db:"id"`
	DivinationID uuid.UUID `db:"divination_id"`
	Name         string    `db:"name"`
	Title        string    `db:"title"`
	URL          string    `db:"url"`
}
