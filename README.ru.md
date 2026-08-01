# Megan Messenger API

[Read in English](README.md)

Корпоративный мессенджер — backend API на Go. Фронтенд — в `frontend/`.

## Технологии

- **Backend**: Go, chi
- **База данных**: PostgreSQL, Redis
- **Связь**: WebSocket (Redis pub/sub)
- **Инструменты**: Docker, SQLC, Swagger, goose

## Auth (как в Telegram)

1. `POST /auth/phone/start` — отправка OTP (в dev код пишется в логи)
2. `POST /auth/phone/verify` — проверка кода → токены **или** challenge пароля
3. Опционально: `POST /auth/password/verify` — если задан cloud password
4. `POST /auth/onboarding/username` — обязателен до доступа к чатам

## Требования

- [Docker](https://www.docker.com/) & Docker Compose
- [Go](https://go.dev/) 1.25+

## Быстрый старт

```bash
cp .env.example .env
make up
make migrate-up
```

Если миграции переписаны с нуля — сбрось volume:

```bash
docker compose down -v
make up
make migrate-up
```

API: http://localhost:8080  
Swagger: http://localhost:8080/docs/index.html

OTP-коды смотри в логах API (`make logs`).

## Структура проекта

```
.
├── cmd/
│   ├── api/            # HTTP + WebSocket сервер
│   └── migrate/        # goose миграции
├── frontend/           # Vue клиент
├── internal/
│   ├── auth/           # phone OTP + JWT
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
