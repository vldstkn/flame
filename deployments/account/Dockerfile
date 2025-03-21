FROM golang:alpine AS builder

WORKDIR /app


RUN apk add --no-cache git

# Копируем сначала только файлы, нужные для зависимостей
COPY go.mod go.sum ./
RUN go mod download


RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Копируем остальные файлы проекта
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o account ./cmd/account.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/account .
COPY --from=builder /app/configs ./configs

EXPOSE 8080

CMD ["./account"]