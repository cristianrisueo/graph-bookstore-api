package store

import "errors"

// Book es el modelo de dominio de un libro.
type Book struct {
	ID       int
	Title    string
	Year     int
	AuthorID int
}

// Errores de dominio de libro
var ErrBookNotFound = errors.New("book not found")
