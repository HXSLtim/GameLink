# 三端 MVP 发布下一阶段执行板

> 来源基线：
> - 当前发布编排：`scripts/mvp-three-end-release.mjs`
> - 当前发布工作流：`.github/workflows/mvp-three-end-release.yml`
> - 当前发布报告：`docs/reports/mvp-three-end-release-report.md`

## 1) 目标
- 在保持当前 `PASS` 能力的前提下，把三端发布从“可跑通”升级为“可持续、可审计、可回滚演练”。

## 2) 当前状态快照
- MVP Core 已可执行并通过：工程验收 + Admin type/lint/test/build。
- 可选回归已具备脚本入口：flow guard / withdraw / full acceptance。
- 仍需推进：发布后报告自动校验、回归开关标准化、owner 交接节奏固化。

## 3) 批次执行（下一阶段）

### Batch D：发布硬化（P0）
- [x] 在 `mvp-three-end-release` 工作流中加入发布报告一致性校验。
- [x] 固化“发布前必选开关组合”：`--with-flow-guard --with-withdraw-regression`。
- [x] 产出失败时的最小回滚动作清单（按 API/Admin/App 分开）。

### Batch E：运行稳定性（P1）
- [ ] 增加发布链路耗时基线（step 级）并跟踪波动。
- [ ] 增加产物保留策略审查（报告/日志 retention 与责任人）。
- [ ] 增加 smoke 范围说明（后端健康 + 前端可达 + 关键 API）。

### Batch F：团队协同（P1）
- [ ] 发布日历：冻结窗口、升级窗口、回滚窗口。
- [ ] 交接模板：Commander -> API/Admin/App Owner 的确认回执。
- [ ] 建立“异常处置 15 分钟规则”（失败即止损并回报）。

## 4) Owner 分配建议
- Commander：发布节奏、开关策略、最终放行。
- API Owner：回归脚本健康、后端日志与数据一致性。
- Admin Owner：管理端质量门禁和构建产物。
- App Owner：用户端质量门禁和工程验收结果。
- QA/SRE：报告归档、工作流运行态、回滚流程核验。

## 5) DoD（完成定义）
- [ ] `mvp-three-end-release` 工作流带报告校验步骤并通过。
- [ ] 发布报告与工程报告均为 `PASS` 且失败步骤为 `0`。
- [ ] 下一阶段执行板每个批次均有 owner 与时间窗口。
