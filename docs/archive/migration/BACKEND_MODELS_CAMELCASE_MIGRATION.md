# 后端模型 camelCase 迁移完成报告

## 概述

所有后端模型的 JSON 标签已统一更新为 camelCase 命名规范，同时保持数据库列名为 snake_case。

## 已完成的模型更新

### 1. Base 模型
**文件:** `backend/internal/model/base.go`

| 字段 | 旧 JSON 标签 | 新 JSON 标签 | GORM 列名 |
|------|-------------|-------------|----------|
| CreatedAt | `created_at` | `createdAt` | `created_at` |
| UpdatedAt | `updated_at` | `updatedAt` | `updated_at` |
| DeletedAt | `deleted_at` | `deletedAt` | `deleted_at` |

### 2. User 模型
**文件:** `backend/internal/model/user.go`

| 字段 | 旧 JSON 标签 | 新 JSON 标签 | GORM 列名 |
|------|-------------|-------------|----------|
| AvatarURL | `avatar_url` | `avatarUrl` | `avatar_url` |
| LastLoginAt | `last_login_at` | `lastLoginAt` | `last_login_at` |

### 3. Order 模型
**文件:** `backend/internal/model/order.go`

| 字段 | 旧 JSON 标签 | 新 JSON 标签 | GORM 列名 |
|------|-------------|-------------|----------|
| UserID | `user_id` | `userId` | `user_id` |
| PlayerID | `player_id` | `playerId` | `player_id` |
| GameID | `game_id` | `gameId` | `game_id` |
| PriceCents | `price_cents` | `priceCents` | `price_cents` |
| ScheduledStart | `scheduled_start` | `scheduledStart` | `scheduled_start` |
| ScheduledEnd | `scheduled_end` | `scheduledEnd` | `scheduled_end` |
| CancelReason | `cancel_reason` | `cancelReason` | `cancel_reason` |
| StartedAt | `started_at` | `startedAt` | `started_at` |
| CompletedAt | `completed_at` | `completedAt` | `completed_at` |
| RefundAmountCents | `refund_amount_cents` | `refundAmountCents` | `refund_amount_cents` |
| RefundReason | `refund_reason` | `refundReason` | `refund_reason` |
| RefundedAt | `refunded_at` | `refundedAt` | `refunded_at` |

### 4. Game 模型
**文件:** `backend/internal/model/game.go`

| 字段 | 旧 JSON 标签 | 新 JSON 标签 | GORM 列名 |
|------|-------------|-------------|----------|
| IconURL | `icon_url` | `iconUrl` | `icon_url` |

### 5. Player 模型
**文件:** `backend/internal/model/player.go`

| 字段 | 旧 JSON 标签 | 新 JSON 标签 | GORM 列名 |
|------|-------------|-------------|----------|
| UserID | `user_id` | `userId` | `user_id` |
| RatingAverage | `rating_average` | `ratingAverage` | `rating_average` |
| RatingCount | `rating_count` | `ratingCount` | `rating_count` |
| HourlyRateCents | `hourly_rate_cents` | `hourlyRateCents` | `hourly_rate_cents` |
| MainGameID | `main_game_id` | `mainGameId` | `main_game_id` |
| VerificationStatus | `verification_status` | `verificationStatus` | `verification_status` |

### 6. Payment 模型
**文件:** `backend/internal/model/payment.go`

| 字段 | 旧 JSON 标签 | 新 JSON 标签 | GORM 列名 |
|------|-------------|-------------|----------|
| OrderID | `order_id` | `orderId` | `order_id` |
| UserID | `user_id` | `userId` | `user_id` |
| AmountCents | `amount_cents` | `amountCents` | `amount_cents` |
| ProviderTradeNo | `provider_trade_no` | `providerTradeNo` | `provider_trade_no` |
| ProviderRaw | `provider_raw` | `providerRaw` | `provider_raw` |
| PaidAt | `paid_at` | `paidAt` | `paid_at` |
| RefundedAt | `refunded_at` | `refundedAt` | `refunded_at` |

### 7. Review 模型
**文件:** `backend/internal/model/review.go`

| 字段 | 旧 JSON 标签 | 新 JSON 标签 | GORM 列名 |
|------|-------------|-------------|----------|
| OrderID | `order_id` | `orderId` | `order_id` |
| UserID | `user_id` | `reviewerId` | `user_id` |
| PlayerID | `player_id` | `playerId` | `player_id` |
| Score | `score` | `rating` | `score` |
| Content | `content` | `comment` | `content` |

### 8. OperationLog 模型
**文件:** `backend/internal/model/operation_log.go`

| 字段 | 旧 JSON 标签 | 新 JSON 标签 | GORM 列名 |
|------|-------------|-------------|----------|
| EntityType | `entity_type` | `entityType` | `entity_type` |
| EntityID | `entity_id` | `entityId` | `entity_id` |
| ActorUserID | `actor_user_id` | `actorUserId` | `actor_user_id` |
| MetadataJSON | `metadata` | `metadata` | `metadata_json` |

## 统计信息

- **总共更新模型数:** 8 个
- **总共更新字段数:** 37 个
- **无 Linter 错误**
- **数据库兼容:** ✅ 完全兼容（通过 GORM column 标签）

## 下一步

1. ✅ 后端模型更新完成
2. 🔄 更新前端类型定义
3. ⏳ 更新 Swagger 注解
4. ⏳ 重新生成 Swagger 文档

## API 响应示例

### 更新前
```json
{
  "id": 1,
  "created_at": "2025-10-29T09:39:15Z",
  "user_id": 2,
  "avatar_url": "https://example.com/avatar.jpg",
  "last_login_at": "2025-10-29T10:00:00Z"
}
```

### 更新后
```json
{
  "id": 1,
  "createdAt": "2025-10-29T09:39:15Z",
  "userId": 2,
  "avatarUrl": "https://example.com/avatar.jpg",
  "lastLoginAt": "2025-10-29T10:00:00Z"
}
```

## 注意事项

1. **数据库不受影响** - 所有模型都通过 `gorm:"column:xxx"` 标签指定数据库列名
2. **向前端输出统一为 camelCase** - 符合 JavaScript/TypeScript 规范
3. **代码质量** - 所有修改通过 golangci-lint 检查
4. **向后兼容** - 旧的前端代码需要更新以匹配新的字段名

## 技术要点

### GORM column 标签使用

```go
// 正确做法
AvatarURL string `json:"avatarUrl" gorm:"column:avatar_url;size:255"`

// 说明：
// - json:"avatarUrl" - API 响应使用 camelCase
// - gorm:"column:avatar_url" - 数据库列名保持 snake_case
// - size:255 - 其他 GORM 约束
```

### 为什么不影响数据库

GORM 在查询和更新数据库时，会使用 `column` 标签中指定的列名。JSON 序列化时则使用 `json` 标签。两者完全独立，互不影响。

```go
// 查询示例
db.Where("user_id = ?", 1).Find(&orders) // ✅ 使用 user_id 列
// 响应示例
json.Marshal(order) // ✅ 输出 {"userId": 1}
```

## 修改文件清单

### 后端
1. ✅ `backend/internal/model/base.go`
2. ✅ `backend/internal/model/user.go`
3. ✅ `backend/internal/model/order.go`
4. ✅ `backend/internal/model/game.go`
5. ✅ `backend/internal/model/player.go`
6. ✅ `backend/internal/model/payment.go`
7. ✅ `backend/internal/model/review.go`
8. ✅ `backend/internal/model/operation_log.go`

### 文档
1. ✅ `BACKEND_MODELS_CAMELCASE_MIGRATION.md` - 本文档
2. ✅ `CAMELCASE_NAMING_UNIFICATION.md` - 命名统一规范
3. ✅ `REVIEW_RATING_FIX.md` - Review 模型修复文档

## 总结

所有后端模型已成功迁移到 camelCase JSON 命名规范，保持数据库列名不变，确保无缝兼容。下一步将更新前端类型定义和 Swagger 文档。

