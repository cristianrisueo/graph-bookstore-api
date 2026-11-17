.PHONY: db-up db-down db-logs migrate-up migrate-down seed generate run setup reset

# Levanta el contenedor de la base de datos en segundo plano
db-up:
	docker compose up -d db

# Detiene y elimina los contenedores de docker compose
db-down:
	docker compose down

# Muestra en vivo los logs de la base de datos
db-logs:
	docker compose logs -f db

# Aplica las migraciones a la base de datos
migrate-up:
	docker compose run --rm migrate

# Revierte todas las migraciones aplicadas
migrate-down:
	docker compose run --rm migrate -path=/migrations -database=postgres://bookstore:bookstore@db:5432/bookstore?sslmode=disable down -all

# Carga datos de ejemplo en la base de datos
seed:
	docker compose exec -T db psql -U bookstore -d bookstore < db/seed.sql

# Genera el código Go a partir de los esquemas GraphQL (gqlgen)
generate:
	go tool gqlgen generate

# Arranca la API en local
run:
	go run .

# Levanta la base de datos, aplica migraciones y carga datos de ejemplo
setup: db-up migrate-up seed

# Borra el volumen de la base de datos y vuelve a dejarlo todo listo
reset:
	docker compose down -v
	$(MAKE) setup
