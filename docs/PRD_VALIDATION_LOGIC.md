# PRD 文档校验逻辑（低维护版）

> 版本：v1.0  
> 更新时间：2026-02-26

## 1. 目标

用自动化校验保证 PRD 体系持续满足“低维护成本”目标，避免文档漂移和重复维护。

## 2. 校验分级

### P0（阻断）

任一失败即 CI 失败：

- SSOT 主链路文件存在且可读：
  - `docs/PRD.md`
  - `docs/PRD_GOVERNANCE.md`
  - `docs/PRODUCT_ROADMAP.md`
  - `docs/PROGRESS.md`
- `PRD.md` 必须声明 SSOT 及治理规则引用。
- `PRD_GOVERNANCE.md` 必须包含：单一事实源、更新规则、Owner 机制。
- `PRODUCT_ROADMAP.md` 必须声明“仅里程碑/时间规划”。
- `PROGRESS.md` 必须声明“仅进度/阻塞项”。
- 核心文档内禁止出现 `app-react`。

### P1（告警）

不阻断，但会提示治理风险：

- 核心文档中出现 `小程序/H5`、`uni-app`（允许在历史版本描述中保留）。
- 文档“最后更新”超过阈值（默认 30 天）。

### P2（观察）

- 章节结构偏离推荐模板（Why/Scope/验收/风险）。
- 关键术语前后不一致（例如“用户端 Web(app)”与“客户端”混用）。

## 3. 输出格式

校验脚本输出两类结果：

- 控制台摘要（PASS/FAIL + 问题数）。
- 报告文件：`docs/reports/prd-validation-report.md`。

## 4. 建议接入方式

### 本地

```bash
node scripts/validate-prd-docs.mjs
```

### CI（推荐）

在 PR 流程中加入：

```bash
node scripts/validate-prd-docs.mjs
```

并将报告文件作为 artifact 保存。

## 5. 维护建议

- 当文档策略变更时，先更新 `PRD_GOVERNANCE.md`，再更新脚本规则。
- 历史文档一律放入 `docs/archive/`，避免污染主链路。
