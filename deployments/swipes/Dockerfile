FROM golang:alpine AS builder

WORKDIR /app

# Копируем сначала только файлы, нужные для зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o swipes ./cmd/swipes.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/swipes .
COPY --from=builder /app/configs ./configs

EXPOSE 8080

CMD ["./swipes"]