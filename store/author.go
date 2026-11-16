package store

import "errors"

// Author is the domain model for an author. Its books are resolved separately.
type Author struct {
	ID   int
	Name string
}

var (
	// ErrAuthorNotFound is returned when no author exists with the given ID.
	ErrAuthorNotFound = errors.New("author not found")
	// ErrAuthorHasBooks is returned when deleting an author that still has books.
	ErrAuthorHasBooks = errors.New("author has books")
)
