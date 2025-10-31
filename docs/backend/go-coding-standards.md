# GameLink Go Coding Standards

本指南规定 GameLink 后端 Go 代码的统一风格与实践，用于提升可读性、可维护性与一致性。所有后端仓库与服务均应遵循本文档，并通过 `golangci-lint` 自动化检查。

适用范围：`backend/` 下所有 Go 代码，包括 `cmd/<service>`、`internal/{config,handler,service,repository,model}` 等模块。

## 基础规范

- 版本与工具
  - Go 1.21+（或仓库 `go.mod` 指定版本）
  - 格式化：`gofmt` + `goimports`（保存即格式化）
  - Lint：`golangci-lint`，配置见仓库根 `.golangci.yml`
- 包与模块
  - 包名短小、小写、无下划线，如 `service`、`handler`、`gormrepo`
  - 每个包聚焦单一职责；避免循环依赖
  - `cmd/<service>/main.go` 仅做组装与启动，不写业务逻辑
- 命名
  - 导出标识使用 UpperCamelCase / MixedCaps：`AdminService`、`OpenPostgres`
  - 非导出标识使用 lowerCamelCase：`parseUintParam`、`listCacheTTL`
  - 常量使用驼峰或 SCREAMING_SNAKE_CASE：`defaultPageSize`、`MaxRetries`、`APP_ENV`
  - 避免模糊缩写；领域术语优先（`Player`, `Payment` 等）
- 目录结构（只列关键）
  - `cmd/<service>/main.go`：服务入口
  - `internal/config`：配置装载与默认值
  - `internal/handler`：HTTP 路由、参数绑定、统一响应
  - `internal/service`：业务逻辑、校验、缓存
  - `internal/repository`：接口定义；`internal/repository/gormrepo` 为实现
  - `internal/model`：实体、枚举与 DTO

## 代码风格

- 格式化与导入
  - 必须通过 `gofmt` 与 `goimports`；导入分组：标准库、第三方、本项目（空行分隔）
  - 未使用的导入与标识不得提交
- 错误处理
  - 失败优先返回：`if err != nil { return ... }`
  - 错误包装使用 `%w` 或 `errors.Join`（Go 1.20+）：`fmt.Errorf("open db: %w", err)`
  - 内部错误消息使用英文；对外统一 `{success, code, message, data}` 响应
  - 约定错误：`service.ErrValidation`、`repository.ErrNotFound` 等，用 `errors.Is` 判断
- 上下文与取消
  - 对外部 I/O（DB、缓存、HTTP）使用 `ctx`：`db.WithContext(ctx)`
  - 不在库中擅自派生无限生命周期的 context
- 并发
  - 使用 `sync`、`time` 等原语前先考虑是否可用服务层缓存/限流替代
  - 共享结构使用互斥保护；避免竞态
- 日志
  - 入口与故障点记录关键日志；敏感信息脱敏
  - 生产环境降低冗余（GORM `logger.Warn`）
- 注释
  - 导出类型与函数具备完整句式注释，首句以被注释对象名开头
  - 复杂逻辑补充动机与不变量说明

## HTTP API 层

- RESTful 路由：`/api/v1/<resource>`，自定义动作示例：`POST /orders/{id}/cancel`
- 统一响应：
  ```json
  {"success":true, "code":200, "message":"OK", "data":{...}, "pagination":{...}}
  ```
- JSON 命名使用 snake_case；Gin 绑定错误返回 400，业务校验返回 400/409，未找到返回 404，内部错误 500
- 中间件
  - 管理端路由加鉴权与限流：`AdminAuth()` + `RateLimitAdmin()`
  - 生产环境必须配置 `ADMIN_TOKEN`，未配置则拒绝（503）

## Service 层

- 只包含业务规则、校验、缓存失效；不感知 HTTP 细节
- 校验失败统一返回 `service.ErrValidation`；找不到实体直接透传 `repository.ErrNotFound`
- 缓存策略
  - 简单列表读缓存：`getCachedList(ctx, cache, key, ttl, fetch)`，写操作后按键失效
  - 内存缓存开发可用，生产优先 Redis

## Repository 层

- `internal/repository` 放接口与分页工具；`gormrepo` 放 GORM 实现
- 查询规范
  - 分页参数归一化：`NormalizePage`、`NormalizePageSize`
  - 计数与列表分步查询；排序明确（一般 `created_at DESC`）
  - 未找到统一返回 `repository.ErrNotFound`，由上层转换为 404
- 更新规范
  - 使用受控字段 map 更新，检查 `RowsAffected`

## Model 层

- 枚举类型（如 `OrderStatus`、`Currency`）提供合规校验方法
- GORM 标签最小化且显式：主键、索引、长度、类型
- 公共和对外 JSON 字段保持稳定；破坏性修改需走版本管理

## 配置与安全

- 配置集中 `configs/`（按环境划分），通过 `internal/config` 加载
- 敏感信息来自环境变量或密钥管理；严禁提交 `.env`
- 生产环境默认关闭 Swagger，可通过 `ENABLE_SWAGGER=false` 控制

## 测试

- 单元测试 `_test.go`，表驱动优先：`Test<Service>_<Case>`
- Service 层覆盖业务路径；Repository 层使用内存或临时 SQLite/容器化 Postgres 做集成测试
- 覆盖率建议≥60%，关键模块更高；`go test -race` 在 CI 运行

## 提交与工作流

- 提交信息遵循 Conventional Commits：`feat(order): 增加优惠券校验`
- 预提交建议执行：`go mod tidy`、`golangci-lint run`、`go test ./...`
- Make 目标
  - `make deps`：`go mod tidy`
  - `make lint`：`golangci-lint run --timeout=5m`
  - `make test`：`go test ./...`
  - `make run CMD=user-service`

## 实用模式与示例

错误处理与包装：
```go
if err := repo.Update(ctx, entity); err != nil {
    if errors.Is(err, repository.ErrNotFound) {
        return nil, ErrNotFound
    }
    return nil, fmt.Errorf("update entity: %w", err)
}
```

分页归一化与响应：
```go
items, total, err := repo.ListPaged(ctx, page, pageSize)
if err != nil { return nil, nil, err }
p := buildPagination(page, pageSize, total)
return items, &p, nil
```

Gin 处理器统一返回：
```go
writeJSON(c, http.StatusOK, model.APIResponse[[]model.Game]{
    Success: true, Code: http.StatusOK, Message: "OK", Data: games, Pagination: pagination,
})
```

## 禁止与注意事项

- 禁止在 `main` 中堆叠业务逻辑；仅组装依赖与启动
- 禁止将 HTTP 上下文对象泄漏至仓储层
- 禁止未显式错误处理（忽略 `err`）
- 严禁在日志或响应中输出密码、令牌等敏感信息

---

更新本标准时，请在 PR 中说明变更动机与影响范围，并同步调整 `.golangci.yml` 与相关流水线配置。

