-- Crea la tabla de autores
CREATE TABLE authors (
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Crea la tabla de libros, cada uno ligado a un autor
CREATE TABLE books (
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title      TEXT        NOT NULL,
    year       INT         NOT NULL,
    author_id  BIGINT      NOT NULL REFERENCES authors(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexa author_id para acelerar las búsquedas de libros por autor
CREATE INDEX books_author_id_idx ON books (author_id);
