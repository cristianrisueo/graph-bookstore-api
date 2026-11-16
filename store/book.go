package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Book es el modelo de dominio de un libro.
type Book struct {
	ID       int
	Title    string
	Year     int
	AuthorID int
}

// Errores de dominio de libro
var ErrBookNotFound = errors.New("book not found")

// ListBooks devuelve todos los libros ordenados por id.
func (s *Store) ListBooks(ctx context.Context) ([]*Book, error) {
	// Ejecuta la consulta y obtiene un iterador de filas.
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, year, author_id
		FROM books
		ORDER BY id
	`)

	if err != nil {
		return nil, fmt.Errorf("list books: %w", err)
	}

	defer rows.Close()

	// Itera sobre las filas y construye la lista de libros.
	var books []*Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Year, &b.AuthorID); err != nil {
			return nil, fmt.Errorf("scan book: %w", err)
		}
		books = append(books, &b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list books: %w", err)
	}

	// Devuelve la lista de libros.
	return books, nil
}

// GetBook devuelve el libro con el id indicado, o ErrBookNotFound si no existe.
func (s *Store) GetBook(ctx context.Context, id int) (*Book, error) {
	// Ejecuta la consulta y obtiene una fila.
	var b Book
	err := s.db.QueryRowContext(ctx, `
		SELECT id, title, year, author_id
		FROM books
		WHERE id = $1
	`, id).Scan(&b.ID, &b.Title, &b.Year, &b.AuthorID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrBookNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get book: %w", err)
	}

	// Devuelve el libro.
	return &b, nil
}

// ListBooksByAuthor devuelve los libros del autor indicado, ordenados por id.
func (s *Store) ListBooksByAuthor(ctx context.Context, authorID int) ([]*Book, error) {
	// Ejecuta la consulta y obtiene un iterador de filas.
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, year, author_id
		FROM books
		WHERE author_id = $1
		ORDER BY id
	`, authorID)

	if err != nil {
		return nil, fmt.Errorf("list books by author: %w", err)
	}

	defer rows.Close()

	// Itera sobre las filas y construye la lista de libros.
	var books []*Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Year, &b.AuthorID); err != nil {
			return nil, fmt.Errorf("scan book: %w", err)
		}
		books = append(books, &b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list books by author: %w", err)
	}

	// Devuelve la lista de libros.
	return books, nil
}
