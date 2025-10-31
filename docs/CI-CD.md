# CI/CD 使用说明

本仓库已完善以下 CI/CD 工作流与部署编排：

- 后端 CI：`.github/workflows/backend-ci.yml`
  - 触发：对 `backend/**` 的 `push`/`pull_request`
  - 步骤：`setup-go`（1.22）→ 依赖缓存 → `golangci-lint` → `go vet` → `go test`（race + coverage）→ `go build`

- 镜像构建与推送：`.github/workflows/docker-images.yml`
  - 触发：推送到 `main` 分支、创建 `v*` 标签，支持手动触发
  - 内容：使用 `Dockerfile` 构建前后端镜像，推送到 GHCR（`ghcr.io`），带缓存
  - 镜像命名：`ghcr.io/<owner>/<repo>-backend`、`ghcr.io/<owner>/<repo>-frontend`
  - 镜像标签：
    - `latest`（main 分支）
    - `sha-<短SHA>`（所有推送）
    - `<分支/标签名>`（斜杠会转为 `-`，并统一为小写）

- 可选部署（SSH）：`.github/workflows/deploy.yml`
  - 触发：
    - 自动：`main` 分支 → 部署到 `staging` 环境；`v*` 标签 → 部署到 `production` 环境
    - 手动：`workflow_dispatch` 可选择目标环境与自定义镜像标签
  - 行为：通过 SSH 将 `docker-compose.prod.yml` 与生成的 `gamelink.env` 上传到服务器，并执行 `docker compose pull && up -d`
  - 说明：仅当部署相关 Secrets 已配置时才会执行（否则自动跳过）

## 部署文件

- 生产编排：`docker-compose.prod.yml`
  - 使用镜像而非本地构建，读取外部环境变量（通过 `--env-file gamelink.env`）
  - 环境变量：
    - `IMAGE_REGISTRY`（默认 `ghcr.io`）
    - `IMAGE_REPOSITORY_PREFIX`（例如 `a2778/gamelink`，需为小写）
    - `IMAGE_TAG`（默认 main 为 `latest`，release 使用标签名）
    - 以及后端所需的敏感配置（见下）

## 必需 Secrets（部署）

在 GitHub 仓库的 `Settings → Secrets and variables → Actions` 中新增：

- 服务器连接：
  - `DEPLOY_HOST`：服务器地址
  - `DEPLOY_PORT`：服务器 SSH 端口（可选，默认 22）
  - `DEPLOY_USER`：SSH 用户名
  - `DEPLOY_SSH_KEY`：SSH 私钥（OpenSSH 格式）
  - `DEPLOY_PATH`：部署路径（例如 `/opt/gamelink`）

- 注册表登录（GHCR）：
  - `GHCR_USERNAME`：GHCR 用户名（通常为 GitHub 用户名/组织名）
  - `GHCR_TOKEN`：具有 `write:packages` 权限的 PAT（用于服务器端拉取）

- 后端运行时敏感配置：
  - `DB_DSN`（Postgres 连接串）
  - `JWT_SECRET_KEY`
  - `CRYPTO_SECRET_KEY`
  - `CRYPTO_IV`
  - `SUPER_ADMIN_EMAIL`
  - `SUPER_ADMIN_PASSWORD`
  - （可选）`REDIS_ADDR`、`REDIS_PASSWORD`、`REDIS_DB`

## 服务器上的 env 文件（示例）

工作流会在部署时生成并上传 `gamelink.env`，其内容示例：

```env
APP_ENV=production
ENABLE_SWAGGER=false
DB_TYPE=postgres
DB_DSN=postgres://user:pass@host:5432/dbname?sslmode=disable
JWT_SECRET_KEY=... 
CRYPTO_SECRET_KEY=...
CRYPTO_IV=...
SUPER_ADMIN_EMAIL=...
SUPER_ADMIN_PASSWORD=...
REDIS_ADDR=...
REDIS_PASSWORD=...
REDIS_DB=0

IMAGE_REGISTRY=ghcr.io
IMAGE_REPOSITORY_PREFIX=a2778/gamelink
IMAGE_TAG=latest
```

## 环境审批与保护

- 建议在仓库 `Settings → Environments` 中创建 `staging` 与 `production` 环境：
  - 配置 Required reviewers（需要审批才能部署）
  - 限制可用的 Secrets 范围与分支

## 使用建议

- 开发环境：使用根目录的 `docker-compose.yml`（本地构建）
- 生产环境：使用 `docker-compose.prod.yml`（拉取 GHCR 镜像），通过 `deploy.yml` 自动化部署
- 若需自定义镜像标签（如回滚），可手动触发 `Deploy (Optional)` 工作流并指定 `image_tag`

## 常见问题

- GHCR 推送失败：检查仓库是否允许 `GITHUB_TOKEN` 推送包，或改用 `GHCR_TOKEN`（PAT）在构建工作流登录并推送。
- 服务器端拉取失败：确认服务器已安装 Docker（含 Compose v2）、开放 22 端口，以及 `GHCR_USERNAME/GHCR_TOKEN` 是否有效。
- 前端无法访问 API：`frontend/nginx.conf` 已将 `/api/` 代理到 `http://backend:8080`，确保 Compose 网络与服务名一致。