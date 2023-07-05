package book

import (
	"errors"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("not found")

type Library struct {
	Book Book
}

type Book struct {
	ID     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	Author string    `json:"author"`
}
