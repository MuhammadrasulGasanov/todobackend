# syntax=docker/dockerfile:1

FROM --platform=linux/amd64 golang:1.24-alpine
# Указание архитектуры для совместимости

# Устанавливаем зависимости
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN go build -o server ./cmd/main.go

# Устанавливаем порт
EXPOSE 8080


# Команда по умолчанию
CMD ["./server"]