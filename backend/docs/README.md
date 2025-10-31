# 文档与备份结构说明

本文档说明仓库内文档、API 规范与备份的组织结构，便于维护与查找。

## 目录结构

- `docs/` — 常规说明与规范文档目录。
  - `swagger.json` / `swagger.yaml` — 对外提供的 API 规范（若存在）。
- `internal/handler/swagger/openapi.json` — 代码内维护的 OpenAPI 规范（推荐作为唯一来源）。
- `archive/docs-backup/` — 文档备份归档目录（历史备份文件），包括：
  - `swagger.json.backup`
  - `swagger.yaml.backup`
- `archive/coverage/` — 历史覆盖率报告与临时输出归档目录（之前散落在根目录的 *_coverage* 目录等）。

## 约定与维护

- API 规范推荐以 `internal/handler/swagger/openapi.json` 为主，并在需要时同步生成 `docs/swagger.json` / `docs/swagger.yaml`。
- 临时覆盖率输出与报告请勿提交到根目录，统一生成到本地或移动到 `archive/coverage/` 后再提交。
- 备份文件统一放置在 `archive/docs-backup/`，避免与主文档混淆。

## 忽略规则

仓库根目录的临时覆盖率目录已在 `.gitignore` 中忽略：

```
/*_coverage*/
/*_cov*/
/*_coverage_final*/
```

如需保留历史报告，请将其移动到 `archive/coverage/` 并提交。

## 常见操作

- 更新 API 规范：
  - 编辑 `internal/handler/swagger/openapi.json`（或使用工具生成），必要时导出到 `docs/swagger.json` / `docs/swagger.yaml`。
- 归档备份：
  - 将旧版 Swagger 备份与导出文件移动到 `archive/docs-backup/`。
- 整理覆盖率输出：
  - 本地运行测试生成 `coverage.out`，如需提交历史报告，请归档到 `archive/coverage/`。

## 说明

本次整理的目标是减少根目录与 `docs/` 的噪音文件，使规范与备份分离、覆盖率输出集中管理，便于后续维护与查阅。