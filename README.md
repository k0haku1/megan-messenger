# Megan Messenger API

[Читать на русском](README.ru.md)

Corporate messenger backend API in Go. Frontend lives in `frontend/`.

## Stack

- **Backend**: Go, chi
- **Database**: PostgreSQL, Redis
- **Realtime**: WebSocket + Redis pub/sub
- **Tooling**: Docker, SQLC, Swagger, goose

## Auth (Telegram-like)

1. `POST /auth/phone/start` — send OTP (logged to console in dev)
2. `POST /auth/phone/verify` — verify code → tokens **or** password challenge
3. Optional: `POST /auth/password/verify` — if cloud password is set
4. `POST /auth/onboarding/username` — required before chats

## Requirements

- [Docker](https://www.docker.com/) & Docker Compose
- [Go](https://go.dev/) 1.25+

## Quick start

```bash
cp .env.example .env
make up
make migrate-up
```

If you changed migrations from scratch, reset the DB volume:

```bash
docker compose down -v
make up
make migrate-up
```

API: http://localhost:8080  
Swagger: http://localhost:8080/docs/index.html

OTP codes are printed in API container logs (`make logs`).

## Project layout

```
.
├── cmd/
│   ├── api/            # HTTP + WebSocket server
│   └── migrate/        # goose migrations
├── frontend/           # Vue client
├── internal/
│   ├── auth/           # phone OTP + JWT
│   ├── user/           # users
│   ├── conversation/   # chats (dm, group)
│   ├── message/        # messages
│   ├── ws/             # WebSocket
│   ├── model/          # domain models
│   ├── repository/     # repository interfaces
│   └── db/             # Postgres (sqlc) + Redis
├── migrations/
├── sql/queries/        # SQL for sqlc
└── docs/               # Swagger
```

## Development

```bash
make sqlc    # regenerate sqlc after editing sql/queries
make swag    # regenerate swagger
make logs    # API container logs
```
