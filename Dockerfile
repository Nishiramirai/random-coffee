# Этап 1: сборка статического бинарного файла
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/bot

# Этап 2: минимальный финальный образ
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/bot .
COPY migrations/ ./migrations/
CMD ["./bot"]
