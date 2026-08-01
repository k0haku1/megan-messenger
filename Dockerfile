
FROM golang:1.26-alpine AS dev

WORKDIR /app

RUN apk add --no-cache git ca-certificates
ENV PATH="${PATH}:/go/bin"
RUN go install github.com/air-verse/air@v1.67.3

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN chmod +x scripts/dev-entrypoint.sh

EXPOSE 8080
ENTRYPOINT ["scripts/dev-entrypoint.sh"]

FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git make ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o megan-messenger-api ./cmd/api


FROM alpine:3.18 AS prod

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/megan-messenger-api .

EXPOSE 8080

CMD ["./megan-messenger-api"]
