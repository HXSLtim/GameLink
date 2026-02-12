# GameLink 全链路联调与数据可信度报告（2026-02-12）

## 范围
- 后端重启并触发迁移/seed（`seed_version=v12`）。
- RBAC 新角色（`csLeader` / `csAgent`）落库与权限校验。
- 用户-陪玩师跨角色业务链路冒烟验证。
- 数据一致性 SQL 检查（`api/scripts/data_integrity_check.sql`）。

## 1) 环境与版本
- 后端服务：`http://127.0.0.1:8080`
- PostgreSQL：`gamelink-postgres`（5433）
- Redis：存在端口占用冲突（6380 被其他容器占用），后端当前可连通缓存服务
- 元数据：
  - `migrate_version = 2026-02-12-v4`
  - `seed_version = 2026-02-12-v12`

## 2) RBAC 角色体系验证结果
- 系统角色已存在 8 个：
  - `superAdmin`, `admin`, `finance`, `customerService`, `csLeader`, `csAgent`, `player`, `user`
- 继承关系：
  - `csLeader.parent_id -> csAgent.id`（已生效，`level=1`）
- 演示账号（seed）：
  - `cs.leader@gamelink.com / CsLeader@123` → RBAC: `csLeader,customerService`
  - `cs.agent@gamelink.com / CsAgent@123` → RBAC: `csAgent`
- 权限差异验证：
  - `GET /api/v1/admin/operation-logs`：`csLeader=200`, `csAgent=403`（符合预期）

## 3) 业务链路冒烟（已跑通）
- 账号：
  - 用户：`demo.user@gamelink.com / User@123456`
  - 陪玩师：`pro.player@gamelink.com / Player@123456`
- 流程：
  1. 用户创建订单（`/user/orders`）✅
  2. 用户钱包支付（`/user/payments`, `method=wallet`）✅
  3. 陪玩师接单（`/player/orders/{id}/accept`）✅
  4. 陪玩师完成（`/player/orders/{id}/complete`）✅
  5. 用户提交评价（`/user/reviews`）✅
  6. 用户查询订单详情（`/user/orders/{id}`）✅
- 关键样例：
  - 订单 `#468` 最终状态：`completed`
  - 支付单 `#267` 状态：`paid`
  - 评价单 `#158` 已创建

## 4) 数据可信度检查结果（首次检查）
执行脚本：`api/scripts/data_integrity_check.sql`

- Orphan（孤儿关联）
  - 大部分为 0
  - `player_schedules.player_id -> players.id = 12`（首次检查发现脏数据）

- Business Consistency（业务一致性）
  - `reviews on non-completed orders = 7`
  - `completed orders without paid/refunded payment = 6`
  - 其余关键项为 0

## 5) 结论与下一步（按优先级）
- P0（建议立即）
  - 处理 6380 端口冲突，固定 GameLink Redis 实例来源（避免跨项目共享 Redis）。
  - ✅ 已执行数据修复脚本：`api/scripts/data_integrity_fix.sql`（已先备份再执行）。
- P1
  - 对“订单状态 vs 评价可提交”规则做强校验，避免非 completed 订单写入评价。
  - 对“completed 订单支付一致性”补回归测试。
- P2
  - 增加“客服主管/客服专员”后台菜单可见性回归用例。

## 6) 数据修复执行记录（已完成）
- 备份文件：
  - `api/backups/gamelink-pre-integrity-fix-20260212-154842.sql`（约 1.24 MB）
- 执行脚本：
  - `api/scripts/data_integrity_fix.sql`
- 执行结果摘要：
  - 删除无效评价链路：`reviews=7`（并清理 `review_replies/review_reports/review_appeals` 关联）
  - 回填已完成订单缺失支付：`payments +6`
  - 清理孤儿排班：`player_schedules -12`
  - 修正“取消订单仍 paid”场景：本轮命中 `0`

## 7) 修复后复检（对比）
- 修复前（关键不一致）：
  - `player_schedules.player_id -> players.id = 12`
  - `reviews on non-completed orders = 7`
  - `completed orders without paid/refunded payment = 6`
- 修复后：
  - `player_schedules.player_id -> players.id = 0`
  - `reviews on non-completed orders = 0`
  - `completed orders without paid/refunded payment = 0`
  - `paid payments on canceled orders = 0`

结论：当前 `data_integrity_check.sql` 所覆盖的核心一致性项已全部归零。

## 8) 服务层硬化（新增）
- 目标：防止“未支付订单被完成”导致的数据回归。
- 规则：
  - 用户端完成订单：`/user/orders/{id}/complete` 前必须存在 `paid` 支付记录。
  - 陪玩师端完成订单：`/player/orders/{id}/complete` 前必须存在 `paid` 支付记录。
  - 管理端完成订单：`AdminService.CompleteOrder/UpdateOrder(status=completed)` 前必须存在 `paid` 支付记录。
- 代码位置：
  - `api/internal/service/order/order.go`
    - 新增 `hasPaidPayment(...)`
    - `CompleteOrder(...)`、`CompleteOrderByPlayer(...)` 增加支付前置校验
  - `api/internal/service/admin/order_service.go`
    - 新增 `hasPaidPayment(...)`
    - `UpdateOrder(...)` 在流转到 `completed` 时强校验
    - `CompleteOrder(...)` 增加显式校验
- 回归测试（2026-02-12）：
  - `go test ./internal/service/order -run "CompleteOrder"` ✅
  - `go test ./internal/service/order` ✅
  - `go test ./internal/service/admin` ✅（当前无测试文件，仅编译校验）

## 9) 联调回归脚本（新增）
- 脚本：`api/scripts/run_flow_guard_regression.ps1`
- 用途：自动验证以下关键链路
  - 未支付订单：`player/user complete` 均返回 `400` 且包含 `order must be paid before completion`
  - 未完成订单：`POST /user/reviews` 返回 `400` 且包含 `订单未完成`
  - 已支付订单：`player complete` 成功，随后 `review` 成功
- 执行命令：
  - `powershell -ExecutionPolicy Bypass -File api/scripts/run_flow_guard_regression.ps1`
- 最近一次执行（2026-02-12）：
  - Case A（未支付拦截）订单 `#471`：A1~A5 全部 PASS
  - Case B（支付后完成）订单 `#472`：B1~B4 全部 PASS

## 10) 最新复验（2026-02-12 晚间）
- 变更：
  - `api/scripts/run_cs_permission_regression.ps1` 在创建争议测试单后，新增 `POST /user/payments`（`method=wechat` mock）步骤。
  - 目的：避免回归脚本产生“已完成但无支付记录”脏数据。
- 回归结果：
  - 客服权限回归：`C0~C9` 全部 PASS（订单 `#500/#501`，争议 `#97/#98`）。
  - 流程守卫回归：`A1~A5`、`B1~B4` 全部 PASS（订单 `#502/#503`）。
- 数据一致性复检：
  - 执行 `api/scripts/run_data_integrity.ps1 -Fix` 后，再执行两次 `api/scripts/run_data_integrity.ps1`。
  - 结果：
    - `reviews on non-completed orders = 0`
    - `multi pending payments on same order = 0`
    - `completed orders without paid/refunded payment = 0`
    - `paid payments on canceled orders = 0`
    - Orphan 关联检查项全部 `0`
- 结论：
  - 当前“全链路回归 + 数据可信度”均通过；
  - 后续仅需补齐“新提现申请（余额前置）”自动化场景。
