# Используем Go как build stage
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем только файлы с зависимостями для кэширования
COPY go.mod go.sum ./

# Упрощаем и объединяем go мод команды
RUN go mod tidy

# Копируем исходный код
COPY . .

RUN go mod vendor

# Создаём директорию для бинарников
RUN mkdir -p /app/dist

# Сборка под Linux (по умолчанию)
RUN GOOS=linux GOARCH=amd64 go build -o /app/dist/protocli-linux ./cmd/main.go

# Сборка под Windows
RUN GOOS=windows GOARCH=amd64 go build -o /app/dist/protocli.exe ./cmd/main.go
