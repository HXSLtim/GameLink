# PRD 治理与低维护策略

> 版本：v1.0  
> 最后更新：2026-02-26  
> 适用范围：`docs/PRD.md`、`docs/PRD_COMPREHENSIVE.md`、`docs/PRODUCT_ROADMAP.md`、`docs/PROGRESS.md`

## 1. 目标

本策略用于降低需求文档维护成本，避免一处变更多处重复修改。

## 2. 单一事实源（SSOT）

- `docs/PRD.md`：唯一需求事实源（范围、优先级、验收口径）。
- `docs/PRODUCT_ROADMAP.md`：仅维护“时间与里程碑”，不重复写功能细节。
- `docs/PROGRESS.md`：仅维护“完成状态与阻塞项”，不重复写需求定义。
- `docs/PRD_COMPREHENSIVE.md`：扩展说明文档，默认从 `PRD.md` 引用，不作为主维护入口。

## 3. 更新规则

- 修改需求范围、优先级、验收标准时：只改 `docs/PRD.md`。
- 修改排期或里程碑时：只改 `docs/PRODUCT_ROADMAP.md`。
- 修改实现进度时：只改 `docs/PROGRESS.md`。
- `docs/PRD_COMPREHENSIVE.md` 每次只做引用同步，不新增独立事实。

## 4. 变更门禁（降低回归维护成本）

每次提交涉及产品需求时，至少检查：

- 是否新增了与 `PRD.md` 冲突的描述。
- 是否将“需求定义”误写进 `PROGRESS.md`。
- 是否在 `ROADMAP` 重复复制了 `PRD` 功能明细。

## 5. 结构模板（推荐）

`PRD.md` 建议仅保留：

- 问题与目标（Why）
- 范围与非范围（Scope/Out of Scope）
- 关键用户故事（P0/P1/P2）
- 验收标准（可测试）
- 依赖与风险

## 6. Owner 机制

- PRD Owner：产品负责人（Product-Manager）
- Tech Reviewer：前后端负责人各一名
- 评审周期：双周一次，超 30 天未更新需标记为“待复核”
