# GameLink 前端高可用 Dockerfile
# 支持后端负载均衡

# 构建阶段
FROM node:20-alpine AS builder

ARG NPM_REGISTRY=https://registry.npmmirror.com

WORKDIR /app

# 设置 npm 镜像
RUN npm config set registry ${NPM_REGISTRY}

# 复制依赖文件
COPY package*.json ./

# 安装依赖
RUN npm ci --prefer-offline --no-audit

# 复制源代码
COPY . .

# 构建
RUN npm run build

# 生产阶段
FROM nginx:alpine

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制高可用 Nginx 配置
COPY nginx.ha.conf /etc/nginx/conf.d/default.conf

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://127.0.0.1/health || exit 1

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]
