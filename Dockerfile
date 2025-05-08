# Сборочный этап
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

# Финальный образ
FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache bash curl

# Копируем необходимые файлы
COPY --from=builder /app/app .
COPY wait-for-it.sh .
COPY templates ./templates/  

RUN chmod +x wait-for-it.sh

EXPOSE 8080

CMD ["./app"]