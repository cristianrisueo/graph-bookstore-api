package store

import "errors"

// Book is the domain model for a book. Its author is resolved from AuthorID.
type Book struct {
	ID       int
	Title    string
	Year     int
	AuthorID int
}

// ErrBookNotFound is returned when no book exists with the given ID.
var ErrBookNotFound = errors.New("book not found")
