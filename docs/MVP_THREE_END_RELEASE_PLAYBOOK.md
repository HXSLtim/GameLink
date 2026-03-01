# 三端 MVP 发布作战手册（Team + Gate）

> 目标：把 API / Admin / App 三端发布从“多人口头协作”升级为“可执行、可追溯、可回滚”的标准流程。

## 1. 发布团队编制（推荐）

### 1.1 角色与职责
- **Release Commander（发布指挥）**：统一节奏、冻结窗口、风险决策、放行签字。
- **API Owner（后端负责人）**：后端测试、回归脚本、数据完整性与回滚预案。
- **Admin Owner（管理端负责人）**：管理端类型/静态检查/测试/构建放行。
- **App Owner（用户端负责人）**：用户端质量门禁、关键路径验证与体验回归。
- **QA & SRE（质量与运维）**：验收报告归档、监控告警确认、变更通告。

### 1.2 批次执行机制
- **Batch A（基础门禁）**：工程验收 + 三端编译/测试阻断。
- **Batch B（关键回归）**：订单守卫、提现链路、全服务流程（按开关启用）。
- **Batch C（发布放行）**：部署工作流、健康检查、回滚与通知。

## 2. 一键发布编排入口

新增脚本：`scripts/mvp-three-end-release.mjs`

### 2.1 本地最小发布验收（推荐）
```bash
node scripts/mvp-three-end-release.mjs --install-deps
```

### 2.1.1 严格模式快速自检（不执行实际步骤）
```bash
node scripts/mvp-three-end-release.mjs \
  --strict-release \
  --with-flow-guard \
  --with-withdraw-regression \
  --dry-run
```

### 2.2 带关键回归的发布验收（发布前）
```bash
node scripts/mvp-three-end-release.mjs \
  --install-deps \
  --strict-release \
  --with-flow-guard \
  --with-withdraw-regression \
  --with-full-acceptance \
  --base-url=http://127.0.0.1:8080/api/v1
```

> `--strict-release` 模式下，必须同时提供 `--with-flow-guard` 与 `--with-withdraw-regression`。

## 3. 编排内容（脚本内置顺序）

### 3.1 MVP Core（阻断）
- 执行 `scripts/engineering-acceptance.mjs`（文档校验 + App Gate + API tests）。
- 执行 Admin 四项门禁：`type-check` / `lint` / `test:run` / `build`。

### 3.2 Regression（按开关）
- `api/scripts/run_flow_guard_regression.ps1`
- `api/scripts/run_withdraw_flow_regression.ps1`
- `api/scripts/run_full_service_flow_acceptance.ps1`

## 4. 报告与追溯
- 发布报告输出：`docs/reports/mvp-three-end-release-report.md`
- 工程验收报告：`docs/reports/engineering-acceptance-report.md`
- 报告一致性校验：`node scripts/validate-mvp-release-report.mjs`
- 下一阶段执行板：`docs/MVP_RELEASE_NEXT_PHASE_EXECUTION_BOARD.md`
- 最小回滚清单：`docs/MVP_RELEASE_MINIMUM_ROLLBACK_CHECKLIST.md`

## 5. 与现有工作流关系
- `ci.yml`：持续集成质量门禁。
- `deploy.yml`：环境部署与回滚。
- `full-service-flow-acceptance.yml`：手动触发全服务链路验收。
- **新增建议**：使用 `mvp-three-end-release.yml` 作为三端 MVP 手动发布编排入口。

## 6. 放行标准（MVP）
- `mvp-three-end-release-report.md` 状态为 `PASS`。
- 失败步骤为 `0`。
- 若启用回归开关，所有 Regression 步骤必须通过。
