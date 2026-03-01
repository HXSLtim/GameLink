# syntax=docker/dockerfile:1.4

# --- Build stage ---
FROM node:20-alpine AS build
WORKDIR /app

ARG NPM_REGISTRY=https://registry.npmmirror.com
RUN npm config set registry ${NPM_REGISTRY}

COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci

COPY . .
RUN npm run build

# --- Runtime stage ---
FROM nginx:alpine AS runtime
WORKDIR /usr/share/nginx/html

RUN rm -rf ./*
COPY --from=build /app/dist .
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
