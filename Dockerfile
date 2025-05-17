# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o proxy-man main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/proxy-man .

EXPOSE 8080

CMD ["./proxy-man"]