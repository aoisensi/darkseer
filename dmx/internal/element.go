package internal

import "github.com/google/uuid"

type Element struct {
	ID         uuid.UUID
	Type       string
	Name       string
	Attributes map[string]any
}
