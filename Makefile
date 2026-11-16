.PHONY: db-up db-down migrate-up migrate-down generate run

# Levanta el contenedor de la base de datos en segundo plano
db-up:
	docker compose up -d db

# Detiene y elimina los contenedores de docker compose
db-down:
	docker compose down

# Aplica las migraciones a la base de datos
migrate-up:
	docker compose run --rm migrate

# Revierte todas las migraciones aplicadas
migrate-down:
	docker compose run --rm migrate -path=/migrations -database=postgres://bookstore:bookstore@db:5432/bookstore?sslmode=disable down -all

# Genera el código Go a partir de los esquemas GraphQL (gqlgen)
generate:
	go tool gqlgen generate

# Arranca la API en local
run:
	go run .
