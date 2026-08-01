# Megan Messenger API

[Read in English](README.md)

Корпоративный мессенджер — backend API на Go. Фронтенд — отдельный репозиторий.

## Технологии

- **Backend**: Go, chi
- **База данных**: PostgreSQL, Redis
- **Связь**: WebSocket (Redis pub/sub)
- **Инструменты**: Docker, SQLC, Swagger, goose

## Требования

- [Docker](https://www.docker.com/) & Docker Compose
- [Go](https://go.dev/) 1.25+

## Быстрый старт

```bash
cp .env.example .env
make up
make migrate-up
```

API: http://localhost:8080  
Swagger: http://localhost:8080/docs/index.html

## Структура проекта

```
.
├── cmd/
│   ├── api/            # HTTP + WebSocket сервер
│   └── migrate/        # goose миграции
├── internal/
│   ├── auth/           # JWT auth
│   ├── user/           # пользователи
│   ├── conversation/   # чаты (dm, group)
│   ├── message/        # сообщения
│   ├── ws/             # WebSocket
│   ├── model/          # доменные модели
│   ├── repository/     # интерфейсы репозиториев
│   └── db/             # Postgres (sqlc) + Redis
├── migrations/
├── sql/queries/        # SQL для sqlc
└── docs/               # Swagger
```

## Разработка

```bash
make sqlc    # перегенерировать sqlc после изменения sql/queries
make swag    # перегенерировать swagger
make logs    # логи API контейнера
```

Базовый код перенесён из [lunar](https://github.com/fluffur/lunar) с удалением Discord-специфики (friends, LiveKit voice, frontend).
