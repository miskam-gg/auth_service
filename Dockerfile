FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /auth_service

FROM gcr.io/distroless/base-debian10

COPY --from=builder /auth_service /auth_service

EXPOSE 8080

CMD ["/auth_service"]
