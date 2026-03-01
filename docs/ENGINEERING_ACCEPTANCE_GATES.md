# 工程化验收门禁（Engineering Acceptance Gates）

> 版本：v1.1  
> 更新时间：2026-02-28

## 1. 目标

把分散的质量检查统一为一个可执行入口，实现本地与 CI 一致的工程化验收。

## 2. 统一验收阶段

### Gate A：文档一致性（阻断）

- `npm run validate:prd`（app）
- `npm run validate:bizstory`（app）

目标：PRD、业务逻辑、用户故事文档可持续维护且无关键漂移。

### Gate B：前端工程质量（阻断）

- `npm run type-check`（app）
- `npm run lint`（app）
- `npm run test:run`（app）
- `npm run build`（app）

目标：用户端代码可类型检查、可静态检查、可单元测试、可构建。

### Gate C：后端可执行性（阻断）

- `go test ./...`（api）

目标：后端测试套件在当前环境可执行。

## 3. 输出与报告

- 控制台输出每个 gate 的 PASS/FAIL。
- 报告文件输出至：`docs/reports/engineering-acceptance-report.md`。

## 4. 一键命令

在 `app` 目录执行：

```bash
npm run accept:engineer
```

## 5. 失败策略

- 任一阻断步骤失败 → 总体 FAIL（非零退出码）。
- 失败详情写入报告，便于 CI 上传 artifact 与排查。
