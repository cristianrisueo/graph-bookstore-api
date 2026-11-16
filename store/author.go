package store

import "errors"

// Author es el modelo de dominio de un autor.
type Author struct {
	ID   int
	Name string
}

// Errores de dominio de autor
var (
	ErrAuthorNotFound = errors.New("author not found")
	ErrAuthorHasBooks = errors.New("author has books")
)
