FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY src/ ./src/
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./src/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8000

CMD ["./main"] 