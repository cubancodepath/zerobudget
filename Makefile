SHELL := /bin/sh

.PHONY: db-up db-down migrate-up migrate-down sqlc-gen api

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

migrate-up:
	goose up

migrate-down:
	goose down

sqlc-gen:
	cd backend && sqlc generate

api:
	set -a; . ./.env; set +a; cd backend && go run ./cmd
