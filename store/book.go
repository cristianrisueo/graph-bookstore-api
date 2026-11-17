package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

// Book es el modelo de dominio de un libro.
type Book struct {
	ID       int
	Title    string
	Year     int
	AuthorID int
}

// BookUpdate contiene los campos a actualizar de un libro. Un campo nil
// significa que no se debe cambiar.
type BookUpdate struct {
	Title    *string
	Year     *int
	AuthorID *int
}

// Errores de dominio de libro
var ErrBookNotFound = errors.New("book not found")

// CreateBook crea un libro con los datos indicados y lo devuelve. Devuelve
// ErrAuthorNotFound si el autor no existe.
func (s *Store) CreateBook(ctx context.Context, title string, year, authorID int) (*Book, error) {
	// Ejecuta el insert y obtiene la fila creada.
	var b Book
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO books (title, year, author_id)
		VALUES ($1, $2, $3)
		RETURNING id, title, year, author_id
	`, title, year, authorID).Scan(&b.ID, &b.Title, &b.Year, &b.AuthorID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgForeignKeyViolation {
		return nil, ErrAuthorNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("create book: %w", err)
	}

	// Devuelve el libro creado.
	return &b, nil
}


// scanBooks consume rows y construye la lista de libros, cerrándolo siempre.
func scanBooks(rows *sql.Rows) ([]*Book, error) {
	defer rows.Close()

	var books []*Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Year, &b.AuthorID); err != nil {
			return nil, fmt.Errorf("scan book: %w", err)
		}
		books = append(books, &b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan book: %w", err)
	}

	return books, nil
}

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

	books, err := scanBooks(rows)
	if err != nil {
		return nil, fmt.Errorf("list books: %w", err)
	}

	// Devuelve la lista de libros.
	return books, nil
}

// ListBooksByAuthorIDs devuelve los libros de los autores indicados, ordenados
// por author_id y luego por id.
func (s *Store) ListBooksByAuthorIDs(ctx context.Context, authorIDs []int) ([]*Book, error) {
	// Ejecuta la consulta y obtiene un iterador de filas.
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, year, author_id
		FROM books
		WHERE author_id = ANY($1)
		ORDER BY author_id, id
	`, authorIDs)

	if err != nil {
		return nil, fmt.Errorf("list books by author ids: %w", err)
	}

	books, err := scanBooks(rows)
	if err != nil {
		return nil, fmt.Errorf("list books by author ids: %w", err)
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

// UpdateBook actualiza los campos indicados del libro con el id dado y lo
// devuelve. Devuelve ErrBookNotFound si no existe, o ErrAuthorNotFound si el
// autor indicado no existe.
func (s *Store) UpdateBook(ctx context.Context, id int, u BookUpdate) (*Book, error) {
	// Ejecuta el update y obtiene la fila actualizada. COALESCE deja cada
	// columna sin cambios cuando el parámetro correspondiente es nil.
	var b Book
	err := s.db.QueryRowContext(ctx, `
		UPDATE books
		SET title = COALESCE($2, title),
			year = COALESCE($3, year),
			author_id = COALESCE($4, author_id)
		WHERE id = $1
		RETURNING id, title, year, author_id
	`, id, u.Title, u.Year, u.AuthorID).Scan(&b.ID, &b.Title, &b.Year, &b.AuthorID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgForeignKeyViolation {
		return nil, ErrAuthorNotFound
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrBookNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update book: %w", err)
	}

	// Devuelve el libro actualizado.
	return &b, nil
}

// DeleteBook elimina el libro con el id indicado. Devuelve ErrBookNotFound si
// no existe.
func (s *Store) DeleteBook(ctx context.Context, id int) error {
	// Ejecuta el delete.
	result, err := s.db.ExecContext(ctx, `
		DELETE FROM books
		WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("delete book: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete book: %w", err)
	}
	if rows == 0 {
		return ErrBookNotFound
	}

	return nil
}
