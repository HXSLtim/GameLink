# syntax=docker/dockerfile:1

# --- Build stage ---
FROM golang:1.22-alpine AS build
WORKDIR /app

# Install git for go modules that require it
RUN apk add --no-cache git ca-certificates && update-ca-certificates

# Cache go modules
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy source and build
COPY . .
# Build a static binary for linux amd64
RUN --mount=type=cache,target=/go/pkg/mod CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/main.go

# --- Runtime stage ---
FROM alpine:3.19 AS runtime
WORKDIR /app

# Copy binary and configs
COPY --from=build /app/server /app/server
COPY configs/ ./configs/

# Prepare runtime data dir (for sqlite dev DB)
RUN mkdir -p /app/var

# Default environment (can be overridden)
ENV APP_ENV=development \
    SERVICE_PORT=8080

EXPOSE 8080

CMD ["/app/server"]