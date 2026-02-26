# 业务逻辑与用户故事文档校验逻辑

> 版本：v1.0  
> 更新时间：2026-02-26

## 1. 校验目标

确保“业务流程定义”和“用户故事定义”在核心文档中持续完整、可追踪、可验证，减少需求文档维护漂移。

## 2. 校验范围（核心）

- `docs/PRD.md`
- `docs/PRD_COMPREHENSIVE.md`
- `docs/PRD_GOVERNANCE.md`

## 3. 规则分级

### P0（阻断）

任一失败即校验失败：

1. 文件存在且可读。
2. `PRD.md` 必须包含“功能架构”章节（`## 二、功能架构`）。
3. `PRD_COMPREHENSIVE.md` 必须包含：
   - `## 四、主要业务流程`
   - `## 七、用户故事与优先级`
4. 用户故事必须至少存在 5 条（`US001` 这类编号）。
5. 用户故事表格行必须包含“验收标准”列，且内容不能为空。
6. 业务流程章节必须包含流程图标记（`mermaid` / `stateDiagram` / `flowchart`）。
7. 核心文档禁止出现过时路径 `app-react`。

### P1（告警）

1. 用户故事编号重复。
2. 用户故事优先级分段缺失（P0/P1/P2 任一缺失）。
3. 缺少 Given/When/Then 风格验收描述（允许为建议项，不阻断）。

### P2（观察）

1. 故事规模过大（单条描述过长）。
2. 业务流程过多但无索引导航。

## 4. 结果输出

- 控制台摘要：PASS/FAIL、阻断数、告警数。
- 报告文件：`docs/reports/business-story-validation-report.md`

## 5. 使用方式

```bash
node scripts/validate-business-story-docs.mjs
```

在 `app` 目录可用：

```bash
npm run validate:bizstory
```
