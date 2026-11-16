.PHONY: db-up db-down migrate-up migrate-down generate run

db-up:
	docker compose up -d db

db-down:
	docker compose down

migrate-up:
	docker compose run --rm migrate

migrate-down:
	docker compose run --rm migrate -path=/migrations -database=postgres://bookstore:bookstore@db:5432/bookstore?sslmode=disable down -all

generate:
	go tool gqlgen generate

run:
	go run ./cmd/api
