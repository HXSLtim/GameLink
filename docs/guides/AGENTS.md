# Repository Guidelines

## 项目结构与模块划分
遵循 `docs/project-structure.md`。后端放在 `backend/`，入口位于 `backend/cmd/<service>/main.go`，通用逻辑在 `backend/internal/{config,handler,service,repository,model}`。前端放在 `frontend/`，部署与 IaC 置于 `deployments/`，脚本维护在 `scripts/`。所有文档写入 `docs/`，架构演进同步更新 `docs/architecture/`。端到端、集成与性能测试按照文档中的 unit/integration/e2e/performance 层级组织进 `tests/`。公共运维资产保存在 `tools/` 与 `monitoring/`，新增目录时在 `docs/project-structure.md` 追加说明。

## 构建、测试与开发命令
后端统一通过 Makefile：

```bash
cd backend && make deps        # go mod tidy 与 vendor 校验
cd backend && make lint        # golangci-lint 执行全量规则
cd backend && make test        # go test ./...（CI 自动附加 -race）
cd backend && make run CMD=order-service  # 本地启动指定服务
cd backend && make build       # 构建全部服务产物
```

前端初始化后补充 `npm run dev|build|lint|test`，并在 `.github/workflows/` 中使用同名命令驱动流水线。常用调试脚本归档到 `scripts/` 并在 README 中注明入口。

## 编码风格与命名约定
严格执行 `docs/go-coding-standards.md`：使用 `gofmt` 与 `goimports`，包名短小且小写，导出类型与函数使用 UpperCamelCase / MixedCaps，常量按驼峰或 SCREAMING_SNAKE_CASE。避免模糊缩写，优先依赖注入与单一职责。预提交阶段启用 `golangci-lint`、`go test`、`go mod tidy` 钩子，复杂模块需补充 `README.md` 解释依赖。前端加入 ESLint + Prettier 并保持与后端一致的语义命名。

## API 设计要点
实现接口时参考 `docs/api-design-standards.md`：坚持 RESTful URL 设计（`/api/v1/<resource>`），动词放在 HTTP 方法，必要时通过 `POST /orders/{id}/cancel` 等自定义动作。统一响应结构 `{success, code, message, data}`，错误遵循标准错误码。分页使用 `page` + `page_size`，支持 `fields`、`sort`、`filter` 查询。所有外部接口默认 HTTPS、鉴权与限流，变更需同步更新 OpenAPI 描述及 Postman 集合。WebSocket 事件字段遵循 `type` 和 `data` 模式，GraphQL Schema 需与 REST 版本保持一致能力。

## 测试准则
Go 测试文件命名为 `_test.go`，函数形如 `Test<Service>_<Case>` 并优先使用表驱动。单测覆盖核心逻辑，集成与 e2e 按文档目录落盘。提交前运行 `go test ./...` 与覆盖率输出 `go test -coverprofile=coverage.out`，必要时补充 `tests/fixtures/` 数据。API 变更需附带 Postman 或自动化脚本更新，性能与安全回归依据 `docs/api-design-standards.md` 的限流与防护章节执行。

## 提交与合并请求
遵循 Conventional Commits（如 `feat(order): 增加优惠券校验`），一次提交聚焦单一需求。PR 必须包含背景、实现概要、测试记录以及 API/UI 变更的截图或示例请求，关联对应 issue，并等待所有必需的工作流通过后再合并。涉及多服务改动时同步列出影响范围和回滚策略。

## 配置与安全
配置集中在 `configs/`，敏感信息通过环境变量或秘密管理服务注入，严禁提交 `.env`。新增变量同步写入 `docs/deployment/` 或 `docs/development/getting-started.md`。发布前执行 `go list -m -u all` 审核依赖，并对监控、日志、Tracing 的更新在 `monitoring/` 与 `docs/operations/` 留档。API 层引入新权限时补充安全测试与限流策略，密钥轮换计划至少每季度复盘一次。所有 环境 变更 需 在 变更 日志 中 记录 以 便 审计 追溯。
