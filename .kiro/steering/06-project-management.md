# 项目管理规则

> AI 辅助项目管理指南，用于维护模块状态和 steering 文档
> 最后更新：2024-12-25

## 模块完整度检查

当用户要求检查模块完整度时，AI 应该：

1. **扫描代码目录**：
   - `backend/internal/model/` - 检查 Model 层
   - `backend/internal/repository/` - 检查 Repository 层
   - `backend/internal/service/` - 检查 Service 层
   - `backend/internal/handler/` - 检查 Handler 层

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
| 核心模块 | 18 | 1 | 0 | 19 |
| 新增业务模块 | 0 | 0 | 3 | 3 |
| 营销模块 | 0 | 1 | 5 | 6 |
| **总计** | **18** | **2** | **8** | **28** |

**整体完整度：约 64%**（18 完成 / 28 总计）

---

### 核心模块状态（已完成）

| 模块 | Model | Repo | Service | Handler | 状态 | 说明 |
|------|-------|------|---------|---------|------|------|
| user | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 用户管理 |
| auth | - | - | ✅ | ✅ | ✅ 完成 | 认证授权 |
| role | ✅ | - | ✅ | ✅ | ✅ 完成 | 角色管理 |
| permission | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 权限管理 |
| menu | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 菜单管理 |
| player | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 陪玩师管理 |
| game | ✅ | ✅ | - | ✅ | ✅ 完成 | 游戏管理 |
| service-item | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 服务项目 |
| order | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 订单管理 |
| payment | ✅ | ✅ | ✅ | - | ✅ 完成 | 支付（内部调用） |
| wallet | ✅ | ✅ | ✅ | - | ✅ 完成 | 钱包（内部调用） |
| withdraw | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 提现管理 |
| review | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 评价管理 |
| dispute | ✅ | ✅ | - | ✅ | ✅ 完成 | 争议处理 |
| chat | ✅ | ✅ | ✅ | - | 🔄 进行中 | 聊天（缺Handler） |
| notification | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 通知系统 |
| sensitive-word | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 敏感词 |
| content | - | ✅ | ✅ | ✅ | ✅ 完成 | 内容管理 |
| statistics | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 统计分析 |

---

### 新增业务模块（仅Model完成）

| 模块 | Model | Repo | Service | Handler | 状态 | 优先级 | 说明 |
|------|-------|------|---------|---------|------|--------|------|
| player-rank | ✅ | ❌ | ❌ | ❌ | 🟡 Model完成 | P1 | 陪玩师等级/认证 |
| order-timeout | ✅ | ❌ | ❌ | ❌ | 🟡 Model完成 | P1 | 订单超时处理 |
| user-block | ✅ | ❌ | ❌ | ❌ | 🟡 Model完成 | P1 | 用户拉黑 |

**下一步建议**：按优先级实现 Repository → Service → Handler

---

### 营销模块状态

| 模块 | Model | Repo | Service | Handler | 状态 | 优先级 | 说明 |
|------|-------|------|---------|---------|------|--------|------|
| vip | ✅ | ❌ | ❌ | ❌ | 🟡 Model完成 | P2 | VIP会员系统 |
| coupon | ✅ | ❌ | ❌ | ❌ | 🟡 Model完成 | P2 | 优惠券系统 |
| recharge | ✅ | ❌ | ❌ | ❌ | 🟡 Model完成 | P2 | 充值系统 |
| activity | ✅ | ❌ | ❌ | ❌ | 🟡 Model完成 | P2 | 活动系统 |
| team | ✅ | - | ✅ | - | 🔄 进行中 | P2 | 团队系统 |
| referral | ✅ | ❌ | ❌ | ❌ | 🟡 Model完成 | P3 | 推荐系统（预留） |

---

### 辅助模块状态

| 模块 | Model | Repo | Service | Handler | 状态 | 说明 |
|------|-------|------|---------|---------|------|------|
| collection-entity | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 收款主体 |
| settlement-company | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 结算公司 |
| routing-rule | - | ✅ | ✅ | ✅ | ✅ 完成 | 路由规则 |
| commission | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 佣金配置 |
| ranking | ✅ | ✅ | ✅ | ✅ | ✅ 完成 | 排行榜 |
| operation-log | ✅ | ✅ | - | - | 🔄 进行中 | 操作日志 |
| user-behavior | ✅ | ✅ | - | ✅ | ✅ 完成 | 用户行为 |

---

## 待实现功能优先级

### P1 - 高优先级（核心业务）

1. **player-rank** - 陪玩师等级/认证系统
   - 缺失：Repository, Service, Handler
   - 依赖：player, game

2. **order-timeout** - 订单超时处理
   - 缺失：Repository, Service, Handler
   - 依赖：order

3. **user-block** - 用户拉黑功能
   - 缺失：Repository, Service, Handler
   - 依赖：user

### P2 - 中优先级（营销功能）

4. **vip** - VIP会员系统
5. **coupon** - 优惠券系统
6. **recharge** - 充值系统
7. **activity** - 活动系统
8. **team** - 团队系统（补充 Handler）

### P3 - 低优先级（预留功能）

9. **referral** - 推荐/邀请系统
10. **chat** - 聊天 Handler

---

## 使用方式

用户可以通过自然语言请求：

- "检查一下项目模块完整度"
- "更新一下 steering 文档"
- "player-rank 模块的 Repository 完成了，帮我更新状态"
- "刚完成了用户拉黑功能，更新一下进度"

AI 会根据请求自动扫描代码、更新相关文档。
