FROM golang:1.24.6 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./main.go

FROM debian:bookworm-slim
WORKDIR /app

# 1. Create a non-root user and group, ensuring the UID/GID is manageable.
# Here we stick to 'appuser' for simplicity.
RUN groupadd -r appuser && useradd -r -g appuser appuser

# 2. Create the images directory
RUN mkdir -p /app/images

# 3. Change ownership of the images directory to the appuser.
# AND give the group WRITE permission (g+w) and others read access (o+r)
RUN chown appuser:appuser /app/images
RUN chmod g+w /app/images
RUN chmod o+r /app/images

COPY --from=builder /app/main .
COPY build/config ./config

EXPOSE 8080

# 4. Switch the process execution to the non-root user
USER appuser
CMD ["./main"]
