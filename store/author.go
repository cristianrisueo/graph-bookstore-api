package store

import "errors"

// Author es el modelo de dominio de un autor. Sus libros se resuelven por separado.
type Author struct {
	ID   int
	Name string
}

var (
	// ErrAuthorNotFound se devuelve cuando no existe ningún autor con el ID indicado.
	ErrAuthorNotFound = errors.New("author not found")
	// ErrAuthorHasBooks se devuelve al borrar un autor que todavía tiene libros.
	ErrAuthorHasBooks = errors.New("author has books")
)
