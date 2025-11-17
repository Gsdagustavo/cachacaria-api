FROM golang:1.24.6 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./main.go

FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && update-ca-certificates

# Create non-root user
RUN groupadd -r appuser && useradd -r -g appuser appuser

# Create images folder
RUN mkdir -p /app/images
RUN chown appuser:appuser /app/images
RUN chmod g+w /app/images
RUN chmod o+r /app/images

COPY --from=builder /app/main .
COPY build/config ./config

EXPOSE 8080

USER appuser
CMD ["./main"]
