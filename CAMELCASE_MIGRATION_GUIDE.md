# GameLink camelCase 命名迁移完整指南

## 📋 总览

本文档提供了将 GameLink 项目从 snake_case 迁移到 camelCase 命名规范的完整指南。

## ✅ 已完成的工作

### 后端模型更新（100% 完成）

所有后端模型已更新为 camelCase JSON 标签：

1. ✅ `backend/internal/model/base.go` - Base 模型
2. ✅ `backend/internal/model/user.go` - User 模型
3. ✅ `backend/internal/model/order.go` - Order 模型
4. ✅ `backend/internal/model/game.go` - Game 模型
5. ✅ `backend/internal/model/player.go` - Player 模型
6. ✅ `backend/internal/model/payment.go` - Payment 模型
7. ✅ `backend/internal/model/review.go` - Review 模型
8. ✅ `backend/internal/model/operation_log.go` - OperationLog 模型

## 🔄 待完成的工作

### 1. 更新前端 TypeScript 类型（关键）

所有前端类型定义需要与后端保持一致。

#### 需要更新的文件

**基础类型:**
- `frontend/src/types/user.ts` - User 和 Player 类型

**实体类型:**
- `frontend/src/types/order.ts` - Order 类型
- `frontend/src/types/game.ts` - Game 类型  
- `frontend/src/types/payment.ts` - Payment 类型
- `frontend/src/types/review.ts` - Review 类型（已部分完成）

**其他类型:**
- `frontend/src/types/stats.ts` - Stats 类型（已完成）
- `frontend/src/types/operation.ts` - OperationLog 类型（如果存在）

#### 字段映射参考

| 后端字段（旧） | 后端字段（新） | 前端需要匹配 |
|---------------|--------------|-------------|
| `created_at` | `createdAt` | `createdAt` |
| `updated_at` | `updatedAt` | `updatedAt` |
| `deleted_at` | `deletedAt` | `deletedAt` |
| `user_id` | `userId` | `userId` |
| `player_id` | `playerId` | `playerId` |
| `game_id` | `gameId` | `gameId` |
| `order_id` | `orderId` | `orderId` |
| `avatar_url` | `avatarUrl` | `avatarUrl` |
| `icon_url` | `iconUrl` | `iconUrl` |
| `price_cents` | `priceCents` | `priceCents` |
| `amount_cents` | `amountCents` | `amountCents` |
| `scheduled_start` | `scheduledStart` | `scheduledStart` |
| `scheduled_end` | `scheduledEnd` | `scheduledEnd` |
| `cancel_reason` | `cancelReason` | `cancelReason` |
| `started_at` | `startedAt` | `startedAt` |
| `completed_at` | `completedAt` | `completedAt` |
| `refund_amount_cents` | `refundAmountCents` | `refundAmountCents` |
| `refund_reason` | `refundReason` | `refundReason` |
| `refunded_at` | `refundedAt` | `refundedAt` |
| `paid_at` | `paidAt` | `paidAt` |
| `last_login_at` | `lastLoginAt` | `lastLoginAt` |
| `rating_average` | `ratingAverage` | `ratingAverage` |
| `rating_count` | `ratingCount` | `ratingCount` |
| `hourly_rate_cents` | `hourlyRateCents` | `hourlyRateCents` |
| `main_game_id` | `mainGameId` | `mainGameId` |
| `verification_status` | `verificationStatus` | `verificationStatus` |
| `provider_trade_no` | `providerTradeNo` | `providerTradeNo` |
| `provider_raw` | `providerRaw` | `providerRaw` |
| `entity_type` | `entityType` | `entityType` |
| `entity_id` | `entityId` | `entityId` |
| `actor_user_id` | `actorUserId` | `actorUserId` |

### 2. 更新前端组件（大量）

所有使用这些字段的组件都需要更新。

#### 主要组件列表

**用户模块:**
- `frontend/src/pages/Users/UserList.tsx`
- `frontend/src/pages/Users/UserDetail.tsx`
- `frontend/src/pages/Users/UserFormModal.tsx`

**订单模块:**
- `frontend/src/pages/Orders/OrderList.tsx`
- `frontend/src/pages/Orders/OrderFormModal.tsx`

**游戏模块:**
- `frontend/src/pages/Games/GameList.tsx`
- `frontend/src/pages/Games/GameFormModal.tsx`

**陪玩师模块:**
- `frontend/src/pages/Players/PlayerList.tsx`
- `frontend/src/pages/Players/PlayerFormModal.tsx`

**支付模块:**
- `frontend/src/pages/Payments/PaymentList.tsx`

**评价模块:**
- `frontend/src/pages/Reviews/ReviewList.tsx` - 已更新
- `frontend/src/pages/Reviews/ReviewFormModal.tsx` - 已更新

**仪表盘:**
- `frontend/src/pages/Dashboard/Dashboard.tsx` - 部分更新

### 3. 更新 Swagger 注解

Swagger 注解需要反映新的 camelCase 字段名。

#### 需要更新的文件

**Handler 文件:**
- `backend/internal/admin/user_handler.go`
- `backend/internal/admin/order_handler.go`
- `backend/internal/admin/game_handler.go`
- `backend/internal/admin/player_handler.go`
- `backend/internal/admin/review_handler.go`
- `backend/internal/admin/stats_handler.go`

#### Swagger 注解示例

```go
// 修改前
// @Success 200 {object} map[string]interface{} "data: {\"user_id\": 1, \"created_at\": \"2025-01-01\"}"

// 修改后
// @Success 200 {object} map[string]interface{} "data: {\"userId\": 1, \"createdAt\": \"2025-01-01\"}"
```

### 4. 重新生成 Swagger 文档

```bash
cd backend
swag init -g cmd/user-service/main.go -o docs/swagger
```

## 🔧 执行步骤

### 步骤1：准备工作

1. **备份当前代码:**
```bash
git add .
git commit -m "backup: before camelCase migration"
git branch feature/camelcase-migration
git checkout feature/camelcase-migration
```

2. **确认后端模型已更新:** ✅ 已完成

3. **创建前端类型映射表** - 见上文"字段映射参考"

### 步骤2：更新前端类型定义

按以下顺序更新前端类型文件：

1. `frontend/src/types/user.ts` - BaseEntity 和 User
2. `frontend/src/types/order.ts` - Order
3. `frontend/src/types/game.ts` - Game
4. `frontend/src/types/payment.ts` - Payment
5. `frontend/src/types/review.ts` - Review（已部分完成）

**BaseEntity 示例:**

```typescript
// 修改前
export interface BaseEntity {
  id: number;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

// 修改后
export interface BaseEntity {
  id: number;
  createdAt: string;
  updatedAt: string;
  deletedAt?: string;
}
```

### 步骤3：更新前端组件

使用查找替换功能批量更新组件中的字段引用：

**VSCode 全局替换示例:**
1. 打开查找替换（Ctrl/Cmd + Shift + H）
2. 搜索：`record\.user_id`
3. 替换：`record.userId`
4. 在 `frontend/src/pages` 目录下替换全部

**需要替换的常见模式:**
- `record.user_id` → `record.userId`
- `record.player_id` → `record.playerId`
- `record.game_id` → `record.gameId`
- `record.created_at` → `record.createdAt`
- `record.avatar_url` → `record.avatarUrl`
- `data.price_cents` → `data.priceCents`
- 等等（见"字段映射参考"）

### 步骤4：更新 Swagger 注解

在每个 handler 文件中更新 Swagger 注解：

1. 打开 handler 文件
2. 找到 `@Success` 和 `@Param` 注解
3. 更新示例 JSON 中的字段名为 camelCase
4. 更新参数说明中的字段名

### 步骤5：重新生成 Swagger 文档

```bash
# 确保已安装 swag
go install github.com/swaggo/swag/cmd/swag@latest

# 生成 Swagger 文档
cd backend
swag init -g cmd/user-service/main.go -o docs/swagger

# 验证生成的文档
cat docs/swagger/swagger.json | grep -i "userId"
```

### 步骤6：测试验证

1. **后端测试:**
```bash
cd backend
go test ./...
```

2. **启动后端:**
```bash
cd backend
go run cmd/user-service/main.go
```

3. **验证 API 响应:**
```bash
curl http://localhost:8080/api/v1/admin/users | jq
# 检查返回的字段是否为 camelCase
```

4. **前端测试:**
```bash
cd frontend
npm run dev
```

5. **功能测试清单:**
- [ ] 用户列表显示正常
- [ ] 订单列表显示正常
- [ ] 游戏列表显示正常
- [ ] 陪玩师列表显示正常
- [ ] 支付列表显示正常
- [ ] 评价列表显示正常（已验证）
- [ ] 仪表盘统计数据显示正常（已验证）
- [ ] 创建/编辑/删除功能正常
- [ ] 筛选和搜索功能正常

### 步骤7：代码审查

1. **检查 TypeScript 类型错误:**
```bash
cd frontend
npm run typecheck
```

2. **检查 Lint 错误:**
```bash
cd frontend
npm run lint
```

3. **检查后端 Lint:**
```bash
cd backend
golangci-lint run
```

### 步骤8：文档更新

更新以下文档：
- [ ] API 文档（Swagger）
- [ ] README.md - 添加迁移说明
- [ ] CHANGELOG.md - 记录破坏性变更

## 🚨 注意事项

### 破坏性变更

这是一个**破坏性变更**！旧的前端代码将无法正常工作，直到完全更新。

### 向后兼容方案

如果需要支持旧的 API 客户端，可以考虑：

1. **API 版本控制:**
   - 保留 `/api/v1` 使用 snake_case
   - 新增 `/api/v2` 使用 camelCase

2. **中间件转换:**
   - 创建中间件自动转换 snake_case ↔ camelCase
   - 对性能有一定影响

3. **分阶段迁移:**
   - 先完成后端
   - 再更新前端
   - 最后废弃旧 API

### 常见问题

**Q: 数据库会受影响吗？**
A: 不会！所有模型都通过 `gorm:"column:xxx"` 标签保持数据库列名为 snake_case。

**Q: 需要迁移数据吗？**
A: 不需要！这只是 API 层面的变更。

**Q: 前端需要全部重写吗？**
A: 不需要重写，只需要批量替换字段名即可。

**Q: Swagger 文档会自动更新吗？**
A: 不会，需要手动更新注解后重新生成。

## 📝 检查清单

### 后端
- [x] Base 模型更新
- [x] User 模型更新  
- [x] Order 模型更新
- [x] Game 模型更新
- [x] Player 模型更新
- [x] Payment 模型更新
- [x] Review 模型更新
- [x] OperationLog 模型更新
- [x] 后端 Lint 检查通过
- [ ] Swagger 注解更新
- [ ] Swagger 文档重新生成
- [ ] 后端测试通过

### 前端
- [ ] BaseEntity 类型更新
- [ ] User 类型更新
- [ ] Order 类型更新
- [ ] Game 类型更新
- [ ] Player 类型更新
- [ ] Payment 类型更新
- [x] Review 类型更新（已完成）
- [x] Stats 类型更新（已完成）
- [ ] 用户模块组件更新
- [ ] 订单模块组件更新
- [ ] 游戏模块组件更新
- [ ] 陪玩师模块组件更新
- [ ] 支付模块组件更新
- [x] 评价模块组件更新（已完成）
- [ ] 仪表盘组件更新
- [ ] 前端 TypeScript 检查通过
- [ ] 前端 Lint 检查通过
- [ ] E2E 测试通过

### 文档
- [x] CAMELCASE_NAMING_UNIFICATION.md
- [x] BACKEND_MODELS_CAMELCASE_MIGRATION.md
- [x] REVIEW_RATING_FIX.md
- [x] CAMELCASE_MIGRATION_GUIDE.md（本文档）
- [ ] API 文档更新
- [ ] README.md 更新
- [ ] CHANGELOG.md 更新

## 🎯 当前状态

**进度:** 30% 完成

- ✅ 后端模型完全更新
- ✅ Review 模块前端更新
- ✅ Stats 相关更新
- 🔄 其他前端类型和组件待更新
- ⏳ Swagger 注解和文档待更新

## 📚 相关文档

- [CAMELCASE_NAMING_UNIFICATION.md](./CAMELCASE_NAMING_UNIFICATION.md) - 命名统一规范
- [BACKEND_MODELS_CAMELCASE_MIGRATION.md](./BACKEND_MODELS_CAMELCASE_MIGRATION.md) - 后端模型迁移详情
- [REVIEW_RATING_FIX.md](./REVIEW_RATING_FIX.md) - Review 模块修复
- [DASHBOARD_COMPLETE_FIX.md](./DASHBOARD_COMPLETE_FIX.md) - 仪表盘修复

## 🚀 快速开始（继续工作）

要继续完成迁移，执行以下命令：

```bash
# 1. 更新前端 BaseEntity 类型
code frontend/src/types/user.ts

# 2. 批量替换前端组件中的字段名
# 在 VSCode 中使用全局查找替换（Ctrl/Cmd + Shift + H）

# 3. 更新 Swagger 注解
code backend/internal/admin/*.go

# 4. 重新生成 Swagger 文档
cd backend && swag init -g cmd/user-service/main.go -o docs/swagger

# 5. 测试
cd backend && go test ./...
cd frontend && npm run typecheck && npm run dev
```

## ✨ 预期收益

完成迁移后的收益：

1. **一致性** - 前后端命名完全统一
2. **可维护性** - 符合各语言生态规范
3. **开发效率** - IDE 自动补全更准确
4. **代码质量** - TypeScript 类型检查更有效
5. **文档质量** - Swagger 文档更规范

---

**最后更新:** 2025-10-29
**更新人:** AI Assistant
**版本:** 1.0

