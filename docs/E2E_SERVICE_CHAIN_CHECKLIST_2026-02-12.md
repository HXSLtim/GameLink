# GameLink 全服务链路检查清单（2026-02-12）

## 一、执行目标

- 验证提现链路：`玩家申请 -> 管理端审核 -> 管理端打款 -> 玩家侧状态回写`
- 验证争议链路：`用户发起 -> 管理端分配 -> 处理完成 -> 用户侧可见处理结果`
- 验证客服权限链路：`客服主管/客服专员在后台权限边界是否符合预期`

---

## 二、本次环境与账号

- 服务地址：`http://localhost:8080/api/v1`
- 健康检查：`GET /healthz` -> `200`
- 测试账号：
  - `admin@gamelink.com / Admin123456`
  - `demo.user@gamelink.com / User@123456`
  - `pro.player@gamelink.com / Player@123456`
  - `cs.leader@gamelink.com / CsLeader@123`
  - `cs.agent@gamelink.com / CsAgent@123`

---

## 三、链路执行结果

### 1) 争议链路

- ✅ `POST /user/orders` 创建订单成功（`orderId=500`、`orderId=501`）
- ✅ `POST /user/payments`（`method=wechat` mock）支付成功
- ✅ `POST /admin/orders/:id/confirm` + `POST /admin/orders/:id/start` 成功
- ✅ `POST /user/orders/:id/dispute` 发起争议成功（`disputeId=97`、`disputeId=98`）
- ✅ `csAgent` 无法分配争议：`POST /admin/disputes/:id/assign` 返回 `403`（符合预期）
- ✅ `csAgent` 可处理争议：`POST /admin/disputes/97/resolve` 返回 `200`
- ✅ `csLeader` 可分配争议：`POST /admin/disputes/98/assign` 返回 `200`
- ✅ `csLeader` 可处理争议：`POST /admin/disputes/98/resolve` 返回 `200`

结论：争议链路与客服权限边界已按设计生效。

### 2) 提现链路

- ✅ `GET /player/earnings/summary` 查询成功（本次 `availableBalance=0`）
- ⚠️ 新建提现请求未触发（可用余额不足 10000 分）
- ✅ 复用历史待处理提现单（`withdrawId=129`）
- ✅ `POST /admin/withdraws/129/approve` 审核通过
- ✅ `POST /admin/withdraws/129/complete` 打款完成
- ✅ `GET /player/earnings/withdraw-history` 玩家侧状态为 `completed`

结论：提现审核与打款链路可用；**新提现申请场景受测试数据余额限制**。

### 3) 客服权限链路

- ✅ `csLeader` 访问 `GET /admin/operation-logs` 返回 `200`
- ✅ `csAgent` 访问 `GET /admin/operation-logs` 返回 `403`（符合权限隔离预期）
- ✅ `csAgent` 可处理争议但不可分配
- ✅ `csLeader` 可分配并处理争议

结论：操作日志权限边界与争议权限分层均正常。

---

## 四、待办（按优先级）

### P0

- 将“争议 + 客服权限 + 流程守卫”脚本纳入固定 CI 时段，避免规则回归。

### P1

- 增加“自动构造可提现余额”的测试前置步骤（或种子数据）：
  - 保证至少 1 个陪玩师 `availableBalance >= 10000`
  - 避免每次回归都依赖历史待处理提现单

### P1

- 增加单脚本回归（CI 可跑）：
  - 争议链路 + 客服权限断言 + 提现链路断言 + 数据一致性断言

---

## 五、回归判定标准（给测试同学）

- 争议链路：
  - 用户可成功发起争议；
  - 后台可分配；
  - 处理后用户侧可看到 `resolved/rejected` 终态。
- 提现链路：
  - 玩家能看到 `pending -> approved -> completed` 的状态流转；
  - 管理端审核与完成操作均成功；
  - 玩家侧历史记录与后台状态一致。
- 客服权限：
  - `csAgent` 不可访问高敏日志；
  - `csLeader` 可访问日志；
  - `csAgent` 可处理争议但不可分配；
  - `csLeader` 可分配并处理争议。
