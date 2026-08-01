# Megan Messenger API

[Читать на русском](README.ru.md)

Corporate messenger backend API in Go. Frontend lives in a separate repository.

## Stack

- **Backend**: Go, chi
- **Database**: PostgreSQL, Redis
- **Realtime**: WebSocket + Redis pub/sub
- **Tooling**: Docker, SQLC, Swagger, goose

## Requirements

- [Docker](https://www.docker.com/) & Docker Compose
- [Go](https://go.dev/) 1.25+

## Quick start

```bash
cp .env.example .env
make up
make migrate-up
```

API: http://localhost:8080  
Swagger: http://localhost:8080/docs/index.html

## Project layout

```
.
├── cmd/
│   ├── api/            # HTTP + WebSocket server
│   └── migrate/        # goose migrations
├── internal/
│   ├── auth/           # JWT auth
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

Bootstrapped from [lunar](https://github.com/fluffur/lunar) with Discord-specific parts removed (friends, LiveKit voice, bundled frontend).
