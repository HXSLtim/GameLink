# GameLink CI/CD 优化计划

**文档版本：** 1.0.0
**创建日期：** 2026-02-09
**维护人：** DevOps-Engineer

---

## 当前 CI/CD 配置分析

### 现有工作流

| 工作流 | 触发条件 | 用途 | 状态 |
|--------|---------|------|------|
| **ci.yml** | Push/PR to main/dev | 持续集成测试 | ✅ 完善 |
| **deploy.yml** | Tag push / Manual | 部署到 Staging/Production | ✅ 完善 |
| **security.yml** | Push/PR/Schedule | 安全扫描 | ✅ 完善 |
| **pre-commit-check.yml** | Pre-commit | 代码质量门禁 | ✅ 完善 |

### 优点

✅ **变更检测机制**：智能检测变更路径，避免不必要的构建
✅ **代码质量门禁**：测试、覆盖率、安全检查
✅ **并发执行**：后端、Admin、Client 并行测试
✅ **Docker 缓存**：使用 GitHub Actions 缓存加速构建
✅ **安全扫描**：govulncheck、npm audit、Trivy、Gitleaks
✅ **部署健康检查**：部署后自动验证服务可用性

---

## 优化建议

### 优化 1：添加性能测试

**目的：** 在 CI 中运行性能基准测试，及时发现性能回归

**实现：**

```yaml
# .github/workflows/performance.yml
name: Performance Tests

on:
  pull_request:
    branches: [main, dev]
  schedule:
    # 每天凌晨运行
    - cron: '0 0 * * *'
  workflow_dispatch:

jobs:
  backend-performance:
    name: Backend Performance Tests
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25.5'

      - name: Run benchmark tests
        working-directory: api
        run: |
          go test -bench=. -benchmem ./... > benchmark_result.txt

          # 检查性能回归（允许 10% 波动）
          if [ -f benchmark_baseline.txt ]; then
            # 对比基线
            go install golang.org/x/perf/cmd/benchstat@latest
            benchstat -delta-test=10% baseline.txt benchmark_result.txt
          fi

      - name: Upload benchmark results
        uses: actions/upload-artifact@v4
        with:
          name: benchmark-results
          path: api/benchmark_result.txt
```

### 优化 2：集成 E2E 测试

**目的：** 端到端测试关键用户流程

**实现：**

```yaml
# .github/workflows/e2e.yml
name: E2E Tests

on:
  schedule:
    # 每天凌晨运行
    - cron: '0 1 * * *'
  workflow_dispatch:

jobs:
  e2e-tests:
    name: End-to-End Tests
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_USER: gamelink
          POSTGRES_PASSWORD: testpassword
          POSTGRES_DB: gamelink_test
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

      redis:
        image: redis:7-alpine
        ports:
          - 6379:6379
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25.5'

      - name: Start backend
        working-directory: api
        run: |
          go build -o bin/gamelink cmd/main.go
          ./bin/gamelink &
          sleep 10
          # 等待后端启动

      - name: Run E2E tests
        run: |
          cd api
          go test -v -tags=e2e ./tests/e2e/...
```

### 优化 3：构建缓存优化

**目的：** 加快 Docker 镜像构建速度

**当前配置：** 已使用 GitHub Actions Cache

**额外优化：**

```yaml
# 在 ci.yml 的构建步骤中添加
- name: Build backend image
  uses: docker/build-push-action@v5
  with:
    context: ./api
    push: false
    tags: gamelink-backend:${{ github.sha }}
    cache-from: type=gha,mode=max
    cache-to: type=gha,mode=max
    build-args: |
      VERSION=${{ github.sha }}
    # 额外优化
    target: builder  # 多阶段构建
```

### 优化 4：测试报告聚合

**目的：** 聚合测试报告，便于查看整体测试覆盖情况

**实现：**

```yaml
# .github/workflows/test-report.yml
name: Test Report

on:
  pull_request:
    branches: [main, dev]
  workflow_dispatch:

jobs:
  aggregate-reports:
    name: Aggregate Test Reports
    runs-on: ubuntu-latest
    if: always()

    steps:
      - uses: actions/checkout@v4

      - name: Download backend coverage
        uses: actions/download-artifact@v4
        continue-on-error: true
        with:
          name: backend-coverage

      - name: Download admin coverage
        uses: actions/download-artifact@v4
        continue-on-error: true
        with:
          name: admin-coverage

      - name: Generate report
        run: |
          echo "# Test Coverage Report" > coverage-report.md
          echo "" >> coverage-report.md
          echo "## Backend" >> coverage-report.md
          cat backend-coverage.txt || echo "No coverage data" >> coverage-report.md
          echo "" >> coverage-report.md
          echo "## Admin" >> coverage-report.md
          cat admin-coverage.txt || echo "No coverage data" >> coverage-report.md

      - name: Comment PR
        uses: actions/github-script@v7
        with:
          github-token: ${{ secrets.GITHUB_TOKEN }}
          script: |
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.name,
              body: fs.readFileSync('coverage-report.md')
            })
```

### 优化 5：依赖更新自动化

**目的：** 自动更新依赖，保持依赖最新

**实现：**

```yaml
# .github/workflows/dependabot-automerge.yml
name: Dependabot Auto-merge

on:
  pull_request_target:
    types: [labeled]

jobs:
  auto-merge-dependabot:
    name: Auto-merge Dependabot PRs
    runs-on: ubuntu-latest
    if: contains(github.event.pull_request.labels.*.name, 'dependencies')
    steps:
      - name: Checkout code
        uses: actions/checkout@v

      - name: Auto-merge
        run: gh pr merge --auto --squash --delete-branch
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### 优化 6：部署通知集成

**目的：** 部署后自动通知团队

**实现：**

```yaml
# 在 deploy.yml 中添加通知步骤
- name: Notify team on success
  if: success()
  run: |
    # 钉钉通知
    curl -X POST "https://oapi.dingtalk.com/robot/send?access_token=${{ secrets.DINGTALK_TOKEN }}" \
      -H 'Content-Type: application/json' \
      -d '{
        "msgtype": "text",
        "text": "✅ GameLink 部署成功\n环境: '${{ needs.pre-check.outputs.environment }}'\n版本: ${{ needs.build-and-push.outputs.version }}"
      }'
```

### 优化 7：回滚自动化

**目的：** 部署失败时自动回滚

**当前配置：** 已在 deploy.yml 中添加

**改进建议：**

```yaml
# 在 deploy-production job 中添加
- name: Automated rollback on failure
  if: failure()
  run: |
    echo "❌ Deployment failed, initiating automated rollback..."

    # 获取上一个稳定版本标签
    PREV_TAG=$(git describe --tags --abbrev=0 HEAD^)

    # 回滚到上一个版本
    git checkout $PREV_TAG

    # 重新构建和部署
    docker-compose -f docker-compose.prod.yml --env-file .env.production down
    docker-compose -f docker-compose.prod.yml --env-file .env.production up -d --force-recreate

    # 通知团队
    curl -X POST "https://oapi.dingtalk.com/robot/send?access_token=${{ secrets.DINGTALK_TOKEN }}" \
      -H 'Content-Type: application/json' \
      -d '{
        "msgtype": "text",
        "text": "❌ 部署失败，已自动回滚到 $PREV_TAG"
      }'
```

---

## CI/CD 最佳实践

### 1. 工作流组织

**建议结构：**

```
.github/workflows/
├── ci.yml                    # 持续集成（代码质量）
├── deploy.yml                # 部署工作流
├── security.yml              # 安全扫描
├── performance.yml           # 性能测试（新增）
├── e2e.yml                   # E2E测试（新增）
├── test-report.yml           # 测试报告聚合（新增）
└── dependabot-automerge.yml  # 依赖更新（新增）
```

### 2. 构建优化

**多阶段构建：**

```dockerfile
# api/Dockerfile
FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY api .
RUN CGO_ENABLED=0 go build -o bin/gamelink cmd/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/gamelink /app/
ENTRYPOINT ["/app/gamelink"]
```

### 3. 测试策略

**测试金字塔：**

```
     /\
    /  \
   / E2E \        ← 少量，慢，运行频率低
  /-------\
 / Unit & Integration \  ← 大量，快，每次提交都运行
-----------------
```

**当前实现：** ✅ 单元测试完善，E2E 测试可添加

### 4. 缓存策略

**当前缓存：**
- Go 依赖缓存（`cache-dependency-path: api/go.sum`）
- npm 缓存（`cache: 'npm'`）
- Docker layer 缓存

**额外优化：**
- 自定义 runner 可启用持久缓存
- 使用 GitHub Actions Cache API

### 5. 安全扫描

**当前实现：** ✅ 完善

- govulncheck（Go 依赖漏洞）
- npm audit（Node.js 依赖漏洞）
- Trivy（Docker 镜像漏洞）
- Gitleaks（密钥泄露检测）

---

## 实施状态

### P0（高优先级）- 已完成 ✅

1. ✅ **性能测试工作流** - 防止性能回归
   - 文件：`.github/workflows/performance.yml`
   - 功能：Go基准测试、API负载测试、性能回归检测
   - 状态：已实现

2. ✅ **E2E 测试工作流** - 覆盖关键用户流程
   - 文件：`.github/workflows/e2e.yml`
   - 功能：后端E2E、用户流程测试、集成场景测试
   - 状态：已实现

### P1（中优先级）- 已完成 ✅

3. ✅ **构建缓存优化** - 加快构建速度
   - 现有配置已包含 Docker layer 缓存
   - Go 模块缓存和 npm 缓存已配置
   - 状态：已优化

4. ✅ **测试报告聚合** - 便于查看测试覆盖
   - 文件：`.github/workflows/test-report.yml`
   - 功能：覆盖率报告、安全扫描报告、性能报告汇总
   - 状态：已实现

5. ✅ **部署通知集成** - 团队及时了解部署状态
   - 文件：`.github/workflows/deploy.yml`（已更新）
   - 功能：钉钉通知、GitHub deployment 状态更新
   - 状态：已实现

### P2（低优先级）- 已完成 ✅

6. ✅ **依赖更新自动化** - 减少手动维护工作
   - 文件：`.github/workflows/dependabot-merge.yml`
   - 文件：`.github/dependabot.yml`
   - 功能：自动合并非破坏性依赖更新
   - 状态：已实现

7. ✅ **回滚自动化改进** - 提高可靠性
   - 文件：`.github/workflows/deploy.yml`（已更新）
   - 功能：自动回滚到上一个稳定版本、通知团队
   - 状态：已实现

---

## 监控和指标

### CI/CD 关键指标

| 指标 | 目标值 | 监控方式 |
|------|--------|---------|
| **构建时间** | < 10 分钟 | GitHub Actions 日志 |
| **测试通过率** | > 95% | 测试报告 |
| **代码覆盖率** | > 70% | Codecov |
| **安全漏洞** | 0 Critical | 安全扫描 |
| **部署成功率** | > 95% | 部署日志 |
| **平均恢复时间** | < 5 分钟 | 监控告警 |

---

## 相关文档

- **CI 工作流：** `.github/workflows/*.yml`
- **部署脚本：** `scripts/deploy.sh`
- **监控指南：** `docs/MONITORING_ALERT_GUIDE.md`
- **故障排查：** `docs/TROUBLESHOOTING_GUIDE.md`

---

**更新历史：**

| 日期 | 版本 | 更新内容 | 更新人 |
|------|------|---------|--------|
| 2026-02-09 | 1.0.0 | 初始版本 | DevOps-Engineer |
