# Build stage
FROM golang:1.24.6 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY build/ .

RUN go build -o main ./main.go  # adjust path if needed

# Runtime stage
FROM debian:bookworm-slim
WORKDIR /app

COPY --from=builder /app/main .

# Copy configs (both, to keep image generic)
COPY build/config ./config

EXPOSE 8080
CMD ["./main"]
