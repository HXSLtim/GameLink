# GameLink 多 Agent 工作任务分发

> 目标：将当前阶段需求拆分为三个互不阻塞的工作流，并为负责的 AI Agent 提供清晰的编码规范提示词，确保并行推进、质量可控。

---

## 工作流 A：车队/订单池能力落地
- **范围**  
  1. 订单池抢单接口（`/player/orders/pool`, `/snatch`, `/release`）与队长派单链路。  
  2. Team & Assignment 数据模型（`Team`, `TeamMember`, `TeamOrderAssignment`, `TeamAssignmentMember`, `TeamPayoutPlan`）。  
  3. 队长端 UI + 队员确认流程，对接收益结算。
- **业务要求**  
  - 支持用户下单时指定「单人/车队」模式，后台能根据 `queueType` 自动进入订单池。  
  - 队长抢单后必须在 5 分钟内完成队员分配，否则自动释放。  
  - 每次分配/回收队员都要记录审计日志，并通知相关成员。  
  - 收益分配方案需在订单完成前全部确认，否则走默认平均分。
- **相关 API（README 规划）**  
  - `GET /api/v1/player/orders/pool`、`POST /api/v1/player/orders/{id}/snatch`、`POST /release`（抢单/释放）。  
  - `POST /api/v1/player/teams`、`GET /player/teams`、`POST /teams/{id}/orders/{orderId}/dispatch`（建队与派单）。  
  - `POST /api/v1/player/teams/{teamId}/orders/{orderId}/payout-plan`（收益方案，规划中，需同时更新 README + Swagger）。  
  - README 数据模型中 “Team/Assignment/PayoutPlan” 字段需同步实现状态。
- **产出**  
  - 数据库迁移脚本（新增表/字段、索引、外键）。  
  - Go 数据模型、仓储与 service 层实现，含并发抢单锁。  
  - Swagger/OpenAPI 更新 + Postman/k6 测试脚本。  
  - README「API & 数据模型」同步更新。
- **验收要求**  
  1. k6 并发脚本模拟 100 个队长抢单，冲突率 <0.5%，日志可追踪。  
  2. 单元/集成测试覆盖抢单、派单、队员确认、收益落地等主要路径，覆盖率 ≥85%。  
  3. 队长/队员前端界面演示流程：抢单 → 分配 → 队员确认 → 订单完成 → 收益入账。  
  4. 数据库迁移和 swagger schema 已合并，`make migrate && make test` 全绿。
- **依赖与假设**  
  - 仅依赖现有 Order/Payment 模型；无需等待社区/通知模块。  
  - 风控与推荐（G6）交由其它工作流，此处只提供 `queueType/requiredMembers` 字段。
- **Agent 编码提示词**  
  - 遵循 Go module 导入顺序：标准库 -> 第三方 -> 项目内部。  
  - GORM 模型新增字段附带 `gorm` tag + json tag，必要时使用 `type:json`.  
  - 所有写 API 需加幂等校验 & 日志（`zap.L()`）。  
  - 并发抢单逻辑必须写 race-safe 单测 (`t.Parallel()` + `sync.WaitGroup`)。  
  - Swagger 通过 `swag fmt && swag init` 更新；提交前运行 `make test`.  
  - **重要**：实现前先阅读 README 中与「订单池/车队」相关的 API 与数据模型章节，确保路径/字段保持一致。

---

## 工作流 B：社区/通知/评价维系
- **范围**  
  1. 用户动态（`/user/feeds` 发布/获取）与内容审核流。  
  2. 评价回复 & 站内通知（`NotificationEvent` 实体，`/notifications` API）。  
  3. 前端用户端页面 `/community` & 通知中心入口。
- **业务要求**  
  - 动态仅允许图片 + 文本，单次发布最多 9 张图，最大 10MB。  
  - 所有动态与评价回复必须经过自动审核 + 人工复审通道，违规内容 10 分钟内下架。  
  - 通知中心需要显示未读数、批量已读，并按优先级触发推送（站内信基础版）。  
  - 用户可举报动态/评价，客服后台可查看、处理并反馈结果。
- **相关 API（README 规划）**  
  - `POST /api/v1/user/feeds`, `GET /api/v1/user/feeds`（动态发布/流）。  
  - `POST /api/v1/player/reviews/{reviewId}/reply`、`POST /api/v1/user/reviews/report`（评价回复/举报）。  
  - `GET /api/v1/notifications`, `POST /api/v1/notifications/read`（通知中心）。  
  - README “社区/维系”章节与“Feed/Notification”数据模型需在开发完成后标记为 Implemented。
- **产出**  
  - Feed/Notification 数据模型与迁移。  
  - API 控制器 + service + repository，含审核/举报接口。  
  - 前端 React 页面（列表、发布弹窗、举报/点赞组件）。  
  - 自动化测试：Go handler 测试 + 前端 component 测试。  
  - README API/页面章节更新，描述状态为「In Progress」→「Implemented」。
- **验收要求**  
  1. POST `/user/feeds` 需通过单元+Newman 测试，包含成功/超限/违规等场景。  
  2. Playwright 用例覆盖「发布 → 审核 → 展示 → 举报 → 下架」全链路。  
  3. 通知中心未读数实时同步（WebSocket 或轮询），并有截图/录屏证明。  
  4. README 和 swagger 均新增 `/user/feeds`、`/notifications`、举报/审核接口。
- **依赖与假设**  
  - 可独立于车队/争议开发；仅需 Auth/用户基础数据。  
  - 推送策略初期只做站内信，后续短信由其它工作流接管。
- **Agent 编码提示词**  
  - 前端使用 `React 18 + TypeScript + Ant Design`（若依赖 UI 库）。  
  - API 响应必须统一 `{ success, data, traceId }`，错误走 `pkg/errors`.  
  - 内容审核逻辑抽象为 `service/feed/moderation.go`，便于替换。  
  - 添加 `Newman`/`Playwright` 用例，并在 README 测试章节登记。  
  - 任何可输入文本的 API 需做长度与敏感词校验（复用 `internal/pkg/safety`）。  
  - **重要**：开发前请先查阅 README 中「社区/通知/评价维系」相关 API & 页面规划，严格对齐字段、路由与状态标记。

---

## 工作流 C：售后争议 & 客服指派增强
- **范围**  
  1. `OrderDispute` 实体 + `/user/orders/{id}/dispute`、`/admin/orders/{id}/mediate` API。  
  2. 客服指派工作台（`/admin/assignments`）的数据流与界面。  
  3. 指派日志、SLA 计时与告警 Hook；与风控/推荐接口对接。
- **业务要求**  
  - 用户可在「服务中/完成后 24h」内发起争议，需上传截图证据。  
  - 客服必须在 SLA（默认 30 分钟）内响应，否则系统自动升级告警。  
  - 指派操作需区分来源（系统推荐/人工指定），并允许一键回退。  
  - 争议处理决定（退款/重派/驳回）需同步至订单状态与通知中心。
- **相关 API（README 规划）**  
  - `POST /api/v1/user/orders/{id}/dispute`、`GET /api/v1/admin/orders/{id}/disputes`（争议接口，需在 README API 区和 Swagger 中新增）。  
  - `GET /api/v1/admin/orders/pending-assign`、`GET /admin/orders/{id}/candidates`、`POST /admin/orders/{id}/assign|assign/cancel`、`POST /admin/orders/{id}/mediate`。  
  - README 中的 “客服指派” 与 “订单争议” 状态要随着实现从「方案评审」升级为「Implemented」，并附示例响应。
- **产出**  
  - 数据模型扩展：`OrderDispute`、`AssignmentSource` 字段、操作日志。  
  - Admin 前端页面（列表、详情、时间轴、触发指派/撤销）。  
  - Webhook/告警脚本（Prometheus Rule 或 Loki Query）和文档。  
  - README 功能与 API 状态从「方案评审」更新为「实现中」。
- **验收要求**  
  1. 集成测试模拟「发起争议 → 客服指派 → 裁决 → 退款」完整流程，并记录 traceId。  
  2. Admin UI 截图/录屏展示指派工作台、时间轴、SLA 倒计时以及告警触发。  
  3. Prometheus/Loki 配置文件或脚本提交到 `ops/alerts/`，并在 README 标注。  
  4. Swagger/README 中的指派、争议接口更新为「Implemented」，含示例。
- **依赖与假设**  
  - 依赖现有 Admin API 基础；与车队分配仅通过 `AssignmentSource` 共享数据。  
  - 风控算法（推荐名单）由其它团队提供 mock endpoint。
- **Agent 编码提示词**  
  - Admin 前端遵循现有 `frontend/src/pages/admin` 结构，使用 `React Query` 拉取数据。  
  - 后端 controller 位于 `internal/handler/admin`, service 放 `internal/service/assignment`.  
  - 记录所有更改到 `operation_log`，并附带 `trace_id`.  
  - 编写 `integration` 测试模拟「提交争议 → 客服处理 → 自动退款」全流程。  
  - 提交前运行 `make lint`, `make test`, `npm run lint`, `npm run test`.  
  - **重要**：任何指派/争议逻辑需先核对 README 中对应 API/验收描述，若有差异先更新 README 再实施。

---

如需新增或调整工作流，请在提交 PR 前同步更新本文件，保证所有 Agent 在同一上下文下协作。
