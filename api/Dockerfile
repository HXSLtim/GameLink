# syntax=docker/dockerfile:1.4

# --- Build stage ---
FROM golang:1.25-alpine AS build
WORKDIR /app

# Build-time configuration for Go modules (overridable via build-args)
ARG GOPROXY=https://goproxy.cn,direct
ARG GOPRIVATE=
ARG GONOSUMDB=
ENV GOPROXY=$GOPROXY \
    GOPRIVATE=$GOPRIVATE \
    GONOSUMDB=$GONOSUMDB

# 使用阿里云镜像源加速 apk 安装
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache git ca-certificates && update-ca-certificates

# Cache go modules
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy source and build
COPY . .
# Build a static binary for linux amd64
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/main.go

# --- Runtime stage ---
FROM alpine:3.19 AS runtime
WORKDIR /app

# 使用阿里云镜像源 + 安装运行时依赖
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache wget ca-certificates && update-ca-certificates

# Copy binary and configs
COPY --from=build /app/server /app/server
COPY configs/ ./configs/

# 创建必要的目录
RUN mkdir -p /app/var /app/logs

# 默认使用生产环境（可通过docker-compose或环境变量覆盖）
ENV APP_ENV=production \
    SERVICE_PORT=8080

EXPOSE 8080

CMD ["/app/server"]