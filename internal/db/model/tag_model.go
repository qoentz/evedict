package model

import "github.com/google/uuid"

type Tag struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
