FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

COPY .env .

EXPOSE 8000

CMD ["./main"]

