# Используем официальное изображение Go для сборки приложения
FROM golang:1.23 as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем приложение
RUN go build -o /auth_service

# Используем минимальное изображение для финального контейнера
FROM gcr.io/distroless/base-debian10

# Копируем скомпилированное приложение из builder контейнера
COPY --from=builder /auth_service /auth_service

# Определяем порт, на котором приложение будет работать
EXPOSE 8080

# Устанавливаем команду по умолчанию для запуска приложения
CMD ["/auth_service"]
