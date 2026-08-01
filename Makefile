ifneq (,$(wildcard .env))
    include .env
    export
endif

.PHONY: up down migrate-up swag sqlc logs

up:
	docker compose up -d

migrate-up:
	docker compose exec api go run ./cmd/migrate up

down:
	docker compose down

swag:
	swag fmt -g cmd/api/main.go
	swag init -g cmd/api/main.go

sqlc:
	sqlc generate

logs:
	docker compose logs api -f
