# СТАДИЯ 1: Сборка (builder)
FROM golang:1.21-alpine AS builder

# Устанавливаем необходимые инструменты (например, git, если нужен для go mod)
RUN apk add --no-cache git

WORKDIR /app

# Копируем файлы с зависимостями (для кэширования слоёв)
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем статически слинкованный бинарник
# -ldflags="-s -w" — удаляет отладочную информацию (уменьшает размер)
# CGO_ENABLED=0 — отключаем CGO для статической сборки
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app ./backend/cmd

# СТАДИЯ 2: Финальный образ (минимальный)
FROM alpine:latest

# Устанавливаем доверенные сертификаты (если приложение делает HTTPS-запросы)
#RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из стадии сборки
COPY --from=builder /app/app .
COPY --from=builder /app/frontend ./frontend
# Открываем порт (только документация, реально не пробрасывает)
EXPOSE 8080

# Команда по умолчанию
CMD ["./app"]