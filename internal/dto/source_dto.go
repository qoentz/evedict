package dto

import "github.com/google/uuid"

type Source struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Title string    `json:"title"`
	URL   string    `json:"url"`
}
