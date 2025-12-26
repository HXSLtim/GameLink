# 项目管理规则

> AI 辅助项目管理指南，用于维护模块状态和 steering 文档
> 最后更新：2025-12-26（recharge 模块完成）

## 模块完整度检查

当用户要求检查模块完整度时，AI 应该：

1. **扫描代码目录**：
   - `api/internal/model/` - 检查 Model 层
   - `api/internal/repository/` - 检查 Repository 层
   - `api/internal/service/` - 检查 Service 层
   - `api/internal/handler/` - 检查 Handler 层

2. **对比各层实现**：
   - 有 Model 但缺少 Repository/Service/Handler 的模块
   - 有 Repository 但缺少 Service/Handler 的模块
   - 识别完整实现的模块

3. **生成状态报告**：
   - 列出各模块的实现状态
   - 标注缺失的组件
   - 给出优先级建议

## Steering 文档自动维护

### 何时更新文档

AI 在以下情况下应主动更新 steering 文档：

1. **完成功能实现后**：
   - 更新 `01-product.md` 的"最近完成"列表
   - 更新相关模块的状态

2. **新增数据模型后**：
   - 更新 `04-data-models.md` 或相关子文档
   - 在 `04c-enums-indexes.md` 添加变更日志

3. **修改枚举/索引后**：
   - 更新 `04c-enums-indexes.md`

4. **项目进度变化时**：
   - 更新 `01-product.md` 的进度表格

### 更新规则

**01-product.md 更新规则：**
```markdown
### 最近完成
- ✅ {功能名称}（{简要描述}）  <!-- 新完成的放最前面 -->
```

**04c-enums-indexes.md 变更日志：**
```markdown
| 日期 | 变更内容 |
|------|----------|
| {YYYY-MM-DD} | {变更描述} |  <!-- 新条目放最前面 -->
```

---

## 模块状态追踪

### 整体进度摘要

| 分类 | 完成 | 进行中 | 仅Model | 总计 |
|------|------|--------|---------|------|
| 核心模块 | 19 | 0 | 0 | 19 |
| 新增业务模块 | 4 | 0 | 0 | 4 |
| 营销模块 | 6 | 0 | 0 | 6 |
| 辅助模块 | 7 | 0 | 0 | 7 |
| **总计** | **36** | **0** | **0** | **36** |

**整体完整度：100%**（36 完成 / 36 总计）

---

### 核心模块状态（已完成）

| 模块 | Model | Repo | Service | Handler | 状态 | 说明 |
|------|-------|------|---------|---------|------|------|
| user | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 用户管理（含 batch/tag/loginHistory） |
| auth | - | - | ✅ | ✅ | ✅ 完成 | 认证授权 |
| role | ✅ | - | ✅ | ✅ | ✅ 完成 | 角色管理 |
| permission | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 权限管理（含 auditLog） |
| menu | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 菜单管理 |
| player | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 陪玩师管理（admin+player端） |
| game | ✅ | ✅ | - | ✅ | ✅ 完成 | 游戏管理 |
| service-item | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 服务项目 |
| order | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 订单管理（含 creation/pricing/validation） |
| payment | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 支付（含 provider/refund） |
| wallet | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 钱包（user端有handler） |
| withdraw | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 提现管理（含 routing） |
| review | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 评价管理（含 reply/report/settings/stats） |
| dispute | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 争议处理（含 assignment） |
| chat | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 聊天（user端有handler） |
| notification | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 通知系统 |
| sensitive-word | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 敏感词 |
| content | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 内容管理（含 category/feed/moderation） |
| statistics | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 统计分析（含 analytics/kpi/monitor） |

---

### 新增业务模块

| 模块 | Model | Repo | Service | Handler | 状态 | 优先级 | 说明 |
|------|-------|------|---------|---------|------|--------|------|
| player-rank | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | P1 | 陪玩师等级/认证（含 gamerank/playerrank/playercertification） |
| order-timeout | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | P1 | 订单超时处理（含配置/日志/客服分配） |
| user-block | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | P1 | 用户拉黑（admin+user端） |
| vip | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | P2 | VIP会员系统（admin+user端） |

**下一步建议**：按优先级实现 P2 营销模块（recharge/activity）

---

### 营销模块状态

| 模块 | Model | Repo | Service | Handler | 状态 | 优先级 | 说明 |
|------|-------|------|---------|---------|------|--------|------|
| vip | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | P2 | VIP会员系统（admin+user端） |
| coupon | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | P2 | 优惠券系统（admin+user端） |
| recharge | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | P2 | 充值系统（admin+user端） |
| activity | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | P2 | 活动系统（admin+user端） |
| team | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | P2 | 团队系统（admin+player端） |
| referral | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | P3 | 推荐系统（admin+user端） |

---

### 辅助模块状态

| 模块 | Model | Repo | Service | Handler | 状态 | 说明 |
|------|-------|------|---------|---------|------|------|
| collection-entity | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 收款主体 |
| settlement-company | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 结算公司 |
| routing-rule | - | ✅ | ✅ | ✅ | ✅ 完成 | 路由规则 |
| commission | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 佣金配置 |
| ranking | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 排行榜 |
| operation-log | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 操作日志（admin service 包含） |
| user-behavior | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 用户行为

---

## 待实现功能优先级

### P1 - 高优先级（核心业务）

所有 P1 模块已完成！

### P2 - 中优先级（营销功能）

1. **team** - 团队系统（补充 Repository + Handler）
   - Model 文件：`team.go`
   - Service 已有：`api/internal/service/team/`

### P3 - 低优先级（预留功能）

所有 P3 模块已完成！

---

## 代码统计

### 后端代码分布

| 层级 | 目录数/文件数 | 说明 |
|------|---------------|------|
| Model | 57 文件 | 数据模型定义 |
| Repository | 37 目录 | 数据访问层 |
| Service | 32 目录 | 业务逻辑层 |
| Handler/Admin | 36 文件 | 管理端接口 |
| Handler/User | 10 文件 | 用户端接口 |
| Handler/Player | 7 文件 | 陪玩师端接口 |
| Router | 5 文件 | 路由配置 |
| WebSocket | 4 文件 | 实时通信 |
| Pkg | 12 目录 | 公共工具包 |

---

## 使用方式

用户可以通过自然语言请求：

- "检查一下项目模块完整度"
- "更新一下 steering 文档"
- "player-rank 模块的 Repository 完成了，帮我更新状态"
- "刚完成了用户拉黑功能，更新一下进度"

AI 会根据请求自动扫描代码、更新相关文档。

---

## PRD 功能需求对照

> 根据 `docs/PRD.md` 第4章功能需求，追踪各功能点实现状态
> 详细功能点状态见 `PROGRESS.md` 第7章

### PRD 功能覆盖统计

| 端 | 总功能点 | 已完成 | 进行中 | 未开始 | 完成率 |
|----|----------|--------|--------|--------|--------|
| 用户端 | 17 | 16 | 1 | 0 | 94% |
| 陪玩师端 | 14 | 13 | 0 | 1 | 93% |
| 管理后台 | 24 | 24 | 0 | 0 | 100% |
| **总计** | **55** | **53** | **1** | **1** | **96%** |

### 未完成的 P0/P1 功能点

| 端 | 功能点 | 优先级 | 依赖模块 | 状态 |
|----|--------|--------|----------|------|
| 用户端 | 充值 | P1 | recharge | 🟡 Model完成 |

### 下一步实现建议

1. **P2 实现 recharge 模块**（解锁充值功能）

2. **P2 实现 activity 模块**（营销功能）
