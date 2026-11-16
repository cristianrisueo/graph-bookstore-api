package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

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

// ListAuthors devuelve todos los autores ordenados por id.
func (s *Store) ListAuthors(ctx context.Context) ([]*Author, error) {
	// Ejecuta la consulta y obtiene un iterador de filas.
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name
		FROM authors
		ORDER BY id
	`)

	if err != nil {
		return nil, fmt.Errorf("list authors: %w", err)
	}

	defer rows.Close()

	// Itera sobre las filas y construye la lista de autores.
	var authors []*Author
	for rows.Next() {
		var a Author
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			return nil, fmt.Errorf("scan author: %w", err)
		}
		authors = append(authors, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list authors: %w", err)
	}

	// Devuelve la lista de autores.
	return authors, nil
}

// GetAuthor devuelve el autor con el id indicado, o ErrAuthorNotFound si no existe.
func (s *Store) GetAuthor(ctx context.Context, id int) (*Author, error) {
	// Ejecuta la consulta y obtiene una fila.
	var a Author
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name
		FROM authors
		WHERE id = $1
	`, id).Scan(&a.ID, &a.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAuthorNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get author: %w", err)
	}

	// Devuelve el autor.
	return &a, nil
}
