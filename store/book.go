package store

import "errors"

// Book es el modelo de dominio de un libro. Su autor se resuelve a partir de AuthorID.
type Book struct {
	ID       int
	Title    string
	Year     int
	AuthorID int
}

// ErrBookNotFound se devuelve cuando no existe ningún libro con el ID indicado.
var ErrBookNotFound = errors.New("book not found")
