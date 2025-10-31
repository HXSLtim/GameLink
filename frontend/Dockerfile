# syntax=docker/dockerfile:1

# --- Build stage ---
FROM node:20-alpine AS build
WORKDIR /app

# Install dependencies with cache
COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci

# Copy source and build
COPY . .
RUN npm run build

# --- Runtime stage ---
FROM nginx:1.25-alpine AS runtime

# Copy built assets
COPY --from=build /app/dist /usr/share/nginx/html

# Nginx config for SPA routing
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]