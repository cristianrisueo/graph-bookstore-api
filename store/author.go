package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
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

// CreateAuthor crea un autor con el nombre indicado y lo devuelve.
func (s *Store) CreateAuthor(ctx context.Context, name string) (*Author, error) {
	// Ejecuta el insert y obtiene la fila creada.
	var a Author
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO authors (name)
		VALUES ($1)
		RETURNING id, name
	`, name).Scan(&a.ID, &a.Name)

	if err != nil {
		return nil, fmt.Errorf("create author: %w", err)
	}

	// Devuelve el autor creado.
	return &a, nil
}

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

// UpdateAuthor actualiza el nombre del autor con el id indicado y lo
// devuelve, o ErrAuthorNotFound si no existe.
func (s *Store) UpdateAuthor(ctx context.Context, id int, name string) (*Author, error) {
	// Ejecuta el update y obtiene la fila actualizada.
	var a Author
	err := s.db.QueryRowContext(ctx, `
		UPDATE authors
		SET name = $2
		WHERE id = $1
		RETURNING id, name
	`, id, name).Scan(&a.ID, &a.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAuthorNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update author: %w", err)
	}

	// Devuelve el autor actualizado.
	return &a, nil
}

// DeleteAuthor elimina el autor con el id indicado. Devuelve ErrAuthorNotFound
// si no existe, o ErrAuthorHasBooks si todavía tiene libros asociados.
func (s *Store) DeleteAuthor(ctx context.Context, id int) error {
	// Ejecuta el delete. La clave foránea de books.author_id arbitra si el
	// autor tiene libros, evitando una condición de carrera entre comprobar
	// y borrar.
	result, err := s.db.ExecContext(ctx, `
		DELETE FROM authors
		WHERE id = $1
	`, id)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == pgRestrictViolation || pgErr.Code == pgForeignKeyViolation) {
		return ErrAuthorHasBooks
	}
	if err != nil {
		return fmt.Errorf("delete author: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete author: %w", err)
	}
	if rows == 0 {
		return ErrAuthorNotFound
	}

	return nil
}
