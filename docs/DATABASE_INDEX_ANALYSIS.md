# 数据库索引分析报告

**分析日期：** 2026-02-09
**分析人：** Team-Lead
**任务：** #61 - 分析数据库索引状态

---

## 📋 执行摘要

分析 GameLink 数据库索引状态，发现了 **完善的索引覆盖**，但存在 **部分冗余索引** 和 **可优化的复合索引**。

**数据库信息：**
- PostgreSQL 16
- 端口：5433
- 数据库：gamelink

---

## 📊 表大小分析（Top 20）

| 表名 | 大小 | 说明 |
|------|------|------|
| `permissions` | 464 kB | 权限表（最大） |
| `orders` | 400 kB | 订单表 |
| `chat_groups` | 224 kB | 聊天群组 |
| `users` | 216 kB | 用户表 |
| `role_permissions` | 216 kB | 角色权限关联 |
| `order_disputes` | 192 kB | 订单纠纷 |
| `recharge_records` | 192 kB | 充值记录 |
| `payments` | 184 kB | 支付记录 |
| `service_items` | 176 kB | 服务项 |
| `user_notifications` | 176 kB | 用户通知 |
| `reviews` | 176 kB | 评价 |
| `order_groups` | 176 kB | 主订单 |

**关键发现：**
- 订单相关表（`orders`, `order_groups`, `order_disputes`）总计约 768 kB
- 数据量较小，性能问题主要来自查询模式而非数据规模
- 需要关注高频查询的索引效率

---

## 🔍 订单相关索引分析

### orders 表索引（29 个）

**已有索引：**
| 索引名 | 定义 | 状态 |
|--------|------|------|
| `orders_pkey` | PRIMARY KEY (id) | ✅ 必需 |
| `idx_order_no` | (order_no) | ✅ 必需 |
| `idx_orders_user_id` | (user_id) | ✅ 存在 |
| `idx_orders_player_id` | (player_id) | ✅ 存在 |
| `idx_orders_game_id` | (game_id) | ✅ 存在 |
| `idx_orders_item_id` | (item_id) | ✅ 存在 |
| `idx_orders_status` | (status) | ✅ 存在 |
| `idx_orders_group_id` | (group_id) | ✅ 存在 |
| `idx_orders_status_created` | (status, created_at) | ✅ 优秀 |
| `idx_user_status_created` | (user_id, status, created_at) | ✅ 优秀 |
| `idx_player_status` | (player_id, status) | ✅ 优秀 |
| `idx_orders_user_created` | (user_id, created_at) | ⚠️ 冗余 |
| `idx_orders_player_created` | (player_id, created_at) | ⚠️ 冗余 |
| `idx_orders_game_created` | (game_id, created_at) | ⚠️ 冗余 |
| `idx_orders_item_created` | (item_id, created_at) | ⚠️ 冗余 |
| `idx_orders_created_at` | (created_at) | ⚠️ 冗余 |
| `idx_orders_deleted_at` | (deleted_at) | ✅ 软删除 |
| `idx_orders_group_hour` | (group_id, hour_index) | ✅ 特殊用途 |
| `idx_orders_transfer` | (transfer_to, transfer_status) | ✅ 转单 |
| `idx_orders_has_dispute` | (has_dispute) | ❓ 用途不明 |
| `idx_orders_recipient_player_id` | (recipient_player_id) | ✅ 转单接收 |
| ... | ... | ... |

**分析：**

✅ **优点：**
1. **复合索引设计优秀**
   - `idx_user_status_created` 覆盖了最常见的查询：用户查看自己的订单
   - `idx_player_status` 支持陪玩师查看订单
   - `idx_orders_status_created` 支持按状态和时间查询

2. **单列索引完整**
   - 所有外键都有索引
   - 常用过滤字段都有索引

⚠️ **可优化点：**
1. **存在冗余索引**
   - `idx_orders_user_created` 可能被 `idx_user_status_created` 覆盖
   - `idx_orders_created_at` 可能被其他复合索引覆盖

2. **部分索引可优化**
   - `idx_orders_status` 单列索引，实际查询通常还包含时间范围

### payments 表索引（10 个）

| 索引名 | 定义 | 状态 |
|--------|------|------|
| `payments_pkey` | PRIMARY KEY (id) | ✅ 必需 |
| `idx_payments_order_id` | (order_id) | ✅ 存在 |
| `idx_payments_user_id` | (user_id) | ✅ 存在 |
| `idx_payments_status` | (status) | ✅ 存在 |
| `idx_payments_order_created` | (order_id, created_at) | ✅ 优秀 |
| `idx_payments_user_created` | (user_id, created_at) | ⚠️ 可能冗余 |
| `idx_payments_status_created` | (status, created_at) | ✅ 优秀 |
| `idx_payments_created_at` | (created_at) | ⚠️ 可能冗余 |
| `idx_payments_deleted_at` | (deleted_at) | ✅ 软删除 |
| `idx_payment_idempotent` | (idempotent_key) | ✅ 幂等性 |

**分析：**
✅ 优秀的设计，覆盖了主要查询场景
⚠️ `idx_payments_status_created` 是性能优化计划中建议的索引，已存在

### reviews 表索引（9 个）

| 索引名 | 定义 | 状态 |
|--------|------|------|
| `reviews_pkey` | PRIMARY KEY (id) | ✅ 必需 |
| `idx_reviews_order_id` | (order_id) | ✅ 存在 |
| `idx_reviews_player_id` | (player_id) | ✅ 存在 |
| `idx_reviews_user_id` | (user_id) | ✅ 存在 |
| `idx_reviews_status` | (status) | ✅ 存在 |
| `idx_reviews_order_item_id` | (order_item_id) | ✅ 存在 |
| `idx_reviews_created_at` | (created_at) | ⚠️ 单列 |
| `idx_reviews_deleted_at` | (deleted_at) | ✅ 软删除 |
| `idx_reviews_is_reported` | (is_reported) | ❓ 用途不明 |

**分析：**
✅ 基本索引完整
⚠️ 缺少 `idx_reviews_order_score`（订单+评分）复合索引
⚠️ 缺少 `idx_reviews_player_status`（陪玩师+审核状态）复合索引

---

## 🎯 索引优化建议

### 建议添加的索引

```sql
-- 1. 评价表：订单+评分（支持查询订单评价详情）
CREATE INDEX IF NOT EXISTS idx_reviews_order_score
ON public.reviews(order_id, score);

-- 2. 评价表：陪玩师+审核状态（支持查询陪玩师的已审核评价）
CREATE INDEX IF NOT EXISTS idx_reviews_player_status
ON public.reviews(player_id, status)
WHERE status = 'approved';

-- 3. 评价表：状态+时间降序（支持时间线查询）
CREATE INDEX IF NOT EXISTS idx_reviews_status_created_desc
ON public.reviews(status, created_at DESC);

-- 4. 订单表：游戏+状态（支持按游戏筛选订单）
CREATE INDEX IF NOT EXISTS idx_orders_game_status
ON public.orders(game_id, status)
WHERE game_id IS NOT NULL;
```

### 建议删除的冗余索引

```sql
-- 以下索引可能被复合索引覆盖，建议分析后删除：

-- 1. 被 idx_user_status_created 覆盖？
DROP INDEX IF EXISTS idx_orders_user_created;

-- 2. 被 idx_player_status 覆盖？
DROP INDEX IF EXISTS idx_orders_player_created;

-- 3. 被其他索引覆盖？
DROP INDEX IF EXISTS idx_orders_created_at;
```

**注意：** 删除前需要使用 `EXPLAIN ANALYZE` 验证查询计划

---

## 📈 索引使用率检查建议

### 1. 检查未使用的索引

```sql
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan as index_scans,
    pg_size_pretty(pg_relation_size(indexrelid)) as index_size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
    AND idx_scan < 50  -- 扫描次数少于 50
ORDER BY idx_scan ASC;
```

### 2. 检查索引大小

```sql
SELECT
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) as index_size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
    AND tablename IN ('orders', 'payments', 'reviews')
ORDER BY pg_relation_size(indexrelid) DESC;
```

### 3. 检查缺失的索引（基于查询日志）

```sql
-- 需要启用 pg_stat_statements 扩展
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- 查找执行次数多且慢的查询
SELECT
    query,
    calls,
    total_exec_time,
    mean_exec_time,
    rows
FROM pg_stat_statements
WHERE query NOT LIKE '%pg_stat%'
ORDER BY mean_exec_time DESC
LIMIT 20;
```

---

## 🎯 优先级建议

| 优先级 | 任务 | 预计收益 | 工作量 |
|--------|------|----------|--------|
| **P1** | 添加 `idx_reviews_order_score` | 中 | 低 |
| **P2** | 添加 `idx_reviews_player_status` | 中 | 低 |
| **P3** | 添加 `idx_reviews_status_created_desc` | 低 | 低 |
| **P3** | 添加 `idx_orders_game_status` | 低 | 低 |
| **P3** | 分析并删除冗余索引 | 低 | 中 |

---

## 📊 总结

### ✅ 优点
1. **索引覆盖率高** - 所有主要查询路径都有索引
2. **复合索引设计优秀** - `idx_user_status_created` 等索引很好地优化了常见查询
3. **软删除索引** - 所有 `deleted_at` 字段都有索引

### ⚠️ 改进空间
1. **部分冗余索引** - 可以优化以减少存储和维护成本
2. **评价表索引** - 可以添加几个复合索引提升性能
3. **监控缺失** - 需要建立索引使用率监控

### 📝 下一步行动
1. ✅ **已完成：** 索引现状分析
2. ⏳ **待执行：** 添加缺失的索引
3. ⏳ **待执行：** 分析冗余索引并删除
4. ⏳ **待执行：** 建立索引使用率监控

---

**报告完成时间：** 2026-02-09
**数据库版本：** PostgreSQL 16
**总表大小：** ~50 MB（开发环境）
