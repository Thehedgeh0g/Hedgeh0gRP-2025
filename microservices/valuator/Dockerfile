# Используем образ Go для сборки
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем только файлы с зависимостями для кэширования
COPY go.mod go.sum ./

# Копируем исходный код
COPY . .

# Упрощаем и объединяем go мод команды
RUN go mod tidy && \
    go mod vendor

# Компиляция приложения
RUN go build -o ./bin/main ./cmd/main.go

# Финальный образ с минимальным размером
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем необходимые файлы из билд-образа
COPY --from=builder /app/bin/main ./bin/main
COPY --from=builder /app/vendor ./vendor
COPY pages pages

# Экспортируем порт
EXPOSE 8082

# Запускаем приложение
CMD ["./bin/main"]
