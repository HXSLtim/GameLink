# 生成文件与归档约定

为保持仓库整洁，以下为测试覆盖率、API 文档与构建产物的整理与忽略约定。

## 常见生成物

- 覆盖率与测试输出：
  - `coverage.out`（单次运行）
  - `coverage_*`、`*_coverage*`、`*_cov*`（历史或批次目录）
  - `test_full_run.log`（测试日志）
- API 文档：
  - `internal/handler/swagger/openapi.json`（唯一来源）
  - `docs/swagger.json`、`docs/swagger.yaml`（导出产物）
  - `docs/docs.go`、`docs/docs_test.go`（工具生成产物）
- 构建产物：
  - `user-service`、`*.exe`、`*.out`（二进制/输出文件）

## 存放位置与忽略

- 归档目录：
  - `archive/coverage/` — 历史覆盖率目录与报告统一归档
  - `archive/docs-backup/` — Swagger 等文档备份统一归档
- `.gitignore` 已忽略：
  - `coverage.out`、`/coverage*/`、`/*_coverage*/`、`/*_cov*/`、`/*_coverage_final*/`
  - `docs/swagger.json`、`docs/swagger.yaml`、`docs/docs.go`、`docs/docs_test.go`、`/docs/swagger/*`
  - `user-service`、`*.exe`、`*.out`
  - `test_full_run.log`

## 生成与导出建议

- 覆盖率：
  - `go test -coverprofile=coverage.out ./...`
  - 如需保留历史报告，请将相关目录移动到 `archive/coverage/` 后提交。
- OpenAPI：
  - 以 `internal/handler/swagger/openapi.json` 为唯一规范来源。
  - 导出到 `docs/swagger.json` / `docs/swagger.yaml` 供外部使用，备份则移动至 `archive/docs-backup/`。
- 构建产物：
  - 建议输出到本地临时目录（如 `./bin`），不要提交到仓库。

## 目的

通过统一的归档与忽略策略，减少临时/生成文件对仓库的干扰，确保文档与规范在正确位置维护，同时保留必要的历史记录。