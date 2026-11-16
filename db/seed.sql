-- Vacía las tablas y reinicia los contadores de identidad para poder
-- ejecutar este script tantas veces como se quiera (idempotente).
TRUNCATE books, authors RESTART IDENTITY;

-- Inseta 4 autores y 9 libros. Los autores se insertan primero para poder
-- referenciarlos en la tabla de libros. Arthur C. Clarke se queda sin libros a propósito.
INSERT INTO authors (name) VALUES
    ('Isaac Asimov'),
    ('Ursula K. Le Guin'),
    ('Philip K. Dick'),
    ('Arthur C. Clarke');

INSERT INTO books (title, year, author_id) VALUES
    ('Foundation', 1951, (SELECT id FROM authors WHERE name = 'Isaac Asimov')),
    ('I, Robot', 1950, (SELECT id FROM authors WHERE name = 'Isaac Asimov')),
    ('The Caves of Steel', 1954, (SELECT id FROM authors WHERE name = 'Isaac Asimov')),

    ('A Wizard of Earthsea', 1968, (SELECT id FROM authors WHERE name = 'Ursula K. Le Guin')),
    ('The Left Hand of Darkness', 1969, (SELECT id FROM authors WHERE name = 'Ursula K. Le Guin')),
    ('The Dispossessed', 1974, (SELECT id FROM authors WHERE name = 'Ursula K. Le Guin')),

    ('Do Androids Dream of Electric Sheep?', 1968, (SELECT id FROM authors WHERE name = 'Philip K. Dick')),
    ('Ubik', 1969, (SELECT id FROM authors WHERE name = 'Philip K. Dick')),
    ('A Scanner Darkly', 1977, (SELECT id FROM authors WHERE name = 'Philip K. Dick'));


