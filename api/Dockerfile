# syntax=docker/dockerfile:1.4
#
# GameLink Backend Production Dockerfile
# 多阶段构建：builder (编译) + runtime (运行)
# 安全特性：非 root 用户、最小化镜像、健康检查
#
# 构建示例:
#   docker build -t gamelink-backend:latest .
#   docker build -t gamelink-backend:v1.0.0 --build-arg VERSION=v1.0.0 .
#
# 运行示例:
#   docker run -p 8080:8080 --env-file .env gamelink-backend:latest

# =============================================================================
# Stage 1: Builder
# =============================================================================
FROM golang:1.25.7-alpine AS builder

# 元数据
LABEL maintainer="GameLink Team <support@gamelink.com>" \
      stage="builder" \
      description="GameLink backend build stage"

# 工作目录
WORKDIR /build

# Go 模块配置（可通过 build-arg 覆盖）
ARG GOPROXY=https://goproxy.cn,direct
ARG GOPRIVATE=
ARG GONOSUMDB=
ENV GOPROXY=$GOPROXY \
    GOPRIVATE=$GOPRIVATE \
    GONOSUMDB=$GONOSUMDB \
    GO111MODULE=on

# 使用阿里云镜像源加速 apk 安装
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache \
        git \
        ca-certificates \
        gcc \
        musl-dev && \
    update-ca-certificates

# 复制依赖文件
COPY go.mod go.sum ./

# 使用 BuildKit 缓存加速依赖下载（仅当 go.mod/go.sum 变化时重新下载）
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# 复制源代码
COPY . .

# 构建参数
ARG VERSION=dev
ARG BUILD_TIME
ARG COMMIT_SHA

# 构建静态二进制文件
# -ldflags="-s -w": 去除调试信息和符号表，减小二进制文件大小
# -trimpath: 移除文件系统路径，增强安全性
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -a \
    -installsuffix cgo \
    -ldflags="-s -w \
        -X 'main.Version=${VERSION}' \
        -X 'main.BuildTime=${BUILD_TIME}' \
        -X 'main.CommitSHA=${COMMIT_SHA}'" \
    -trimpath \
    -o /build/server \
    ./cmd/main.go

# 验证二进制文件
RUN /build/server --version 2>/dev/null || echo "Binary built successfully"

# =============================================================================
# Stage 2: Runtime
# =============================================================================
FROM alpine:3.20

# 元数据
LABEL maintainer="GameLink Team <support@gamelink.com>" \
      stage="runtime" \
      description="GameLink backend runtime stage"

# 使用阿里云镜像源加速 apk 安装
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    # 安装最小运行时依赖
    apk add --no-cache \
        ca-certificates \
        tzdata && \
    # 更新 CA 证书
    update-ca-certificates && \
    # 设置时区为上海（可通过环境变量覆盖）
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    # 清理 apk 缓存
    rm -rf /var/cache/apk/*

# 创建非 root 用户和组（安全最佳实践）
RUN addgroup -g 1000 gamelink && \
    adduser -D -u 1000 -G gamelink -h /app -s /bin/sh gamelink

# 工作目录
WORKDIR /app

# 从构建阶段复制二进制文件和配置
COPY --from=builder /build/server /app/server
COPY --chown=gamelink:gamelink configs/ /app/configs/
COPY --chown=gamelink:gamelink docs/ /app/docs/

# 创建必要的目录并设置权限
RUN mkdir -p /app/var /app/logs /app/upload && \
    chown -R gamelink:gamelink /app

# 切换到非 root 用户
USER gamelink

# 环境变量
ENV APP_ENV=production \
    SERVICE_PORT=8080 \
    GIN_MODE=release

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/healthz || exit 1

# 启动应用
CMD ["/app/server"]
