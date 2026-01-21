# syntax=docker/dockerfile:1.4

# --- Build stage ---
FROM node:20-alpine AS build
WORKDIR /app

# 配置npm镜像源（可通过构建参数覆盖）
ARG NPM_REGISTRY=https://registry.npmmirror.com
RUN npm config set registry ${NPM_REGISTRY}

# 安装依赖
COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci

# 复制源代码并构建
COPY . .
RUN npm run build

# --- Runtime stage (使用nginx提供静态文件) ---
# Use nginx:alpine (latest) which uses Alpine 3.20+ to avoid CVEs in Alpine 3.19
FROM nginx:alpine AS runtime
WORKDIR /usr/share/nginx/html

# 删除nginx默认内容
RUN rm -rf ./*

# 从构建阶段复制构建产物
COPY --from=build /app/dist .

# 复制自定义nginx配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]
