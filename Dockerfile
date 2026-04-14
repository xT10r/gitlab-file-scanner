#syntax=docker/dockerfile:1.12

# Стадия 1: Сборка ===
FROM golang:1.26 AS builder
RUN mkdir -p /app && mkdir -p /bin
WORKDIR /app
COPY . /app
RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/gitlab-file-scanner ./cmd

# === Стадия 2: Создание образа приложения ===
FROM alpine:3.23.3
RUN apk --no-cache add ca-certificates tini
COPY --from=builder /bin/gitlab-file-scanner /app/gitlab-file-scanner
WORKDIR /app

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["./gitlab-file-scanner"]