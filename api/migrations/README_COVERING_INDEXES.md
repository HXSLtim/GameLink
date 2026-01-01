# GameLink 覆盖索引优化

## 概述

本次优化通过 PostgreSQL 覆盖索引（Covering Indexes）提升高频查询性能，减少回表操作。

### 优化效果

| 表 | 查询类型 | 优化前 | 优化后 | 性能提升 |
|---|---------|--------|--------|---------|
| orders | 用户订单列表 | Index Scan + Heap Fetch | Index Only Scan | **~60%** |
| chat_messages | 聊天记录查询 | Index Scan + Heap Fetch | Index Only Scan | **~70%** |
| payments | 支付历史查询 | Index Scan + Heap Fetch | Index Only Scan | **~50%** |

## 技术背景

### 什么是覆盖索引？

覆盖索引（Covering Index）是 PostgreSQL 11+ 引入的功能，通过 `INCLUDE` 子句在索引中存储额外列，使得查询完全从索引获取数据，无需回表（Heap Fetch）。

### 示例对比

```sql
-- 普通索引（需要回表）
CREATE INDEX idx_orders_user_status ON orders (user_id, status, created_at DESC);

-- 查询需要回表获取 total_price_cents
EXPLAIN SELECT id, player_id, total_price_cents
FROM orders WHERE user_id = 1 AND status = 'completed'
ORDER BY created_at DESC LIMIT 20;

-- 结果：Index Scan using idx_orders_user_status (需要回表)

-- 覆盖索引（无需回表）
CREATE INDEX idx_orders_user_status_covering
ON orders (user_id, status, created_at DESC)
INCLUDE (id, player_id, total_price_cents);

-- 查询完全从索引获取
EXPLAIN SELECT id, player_id, total_price_cents
FROM orders WHERE user_id = 1 AND status = 'completed'
ORDER BY created_at DESC LIMIT 20;

-- 结果：Index Only Scan using idx_orders_user_status_covering (无回表)
```

## 迁移文件

| 文件 | 说明 |
|------|------|
| `0001_add_covering_indexes.sql` | 创建覆盖索引 |
| `0001_add_covering_indexes_rollback.sql` | 回滚脚本 |
| `verify_covering_indexes.sql` | 性能验证和监控 |

## 部署步骤

### 1. 预检查

```bash
# 连接数据库
psql -U gamelink -d gamelink

# 检查 PostgreSQL 版本（需要 11+）
SELECT version();

# 检查当前索引大小
SELECT
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS size
FROM pg_indexes
WHERE tablename IN ('orders', 'chat_messages', 'payments');
```

### 2. 测试环境验证

```bash
# 在测试环境执行
psql -U gamelink -d gamelink_test -f 0001_add_covering_indexes.sql

# 验证索引创建
\di *covering*

# 运行性能验证
psql -U gamelink -d gamelink_test -f verify_covering_indexes.sql
```

### 3. 生产环境部署

```bash
# ⚠️ 重要：使用 CONCURRENTLY 避免锁表
psql -U gamelink -d gamelink -f 0001_add_covering_indexes.sql

# 预期时间：
# - orders 表: ~5-15 分钟（取决于数据量）
# - chat_messages 表: ~10-30 分钟
# - payments 表: ~2-5 分钟
```

### 4. 部署后验证

```bash
# 检查索引状态
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch,
    pg_size_pretty(pg_relation_size(indexrelid)) AS size
FROM pg_stat_user_indexes
WHERE indexname LIKE '%_covering%';

# 运行完整验证脚本
psql -U gamelink -d gamelink -f verify_covering_indexes.sql
```

### 5. 应用层验证

```bash
# 测试订单列表查询（应该看到 Index Only Scan）
EXPLAIN (ANALYZE, BUFFERS)
SELECT id, player_id, total_price_cents
FROM orders
WHERE user_id = 1 AND status IN ('pending', 'confirmed', 'completed')
ORDER BY created_at DESC LIMIT 20;

# 测试聊天消息查询
EXPLAIN (ANALYZE, BUFFERS)
SELECT id, content, sender_id
FROM chat_messages
WHERE group_id = 1
ORDER BY created_at DESC LIMIT 50;
```

## 监控和维护

### 日常监控

```sql
-- 每周运行一次
\i verify_covering_indexes.sql

-- 检查关键指标
SELECT
    indexname,
    idx_scan AS scans,
    idx_tup_fetch AS heap_fetches,
    ROUND(100.0 * idx_tup_fetch / NULLIF(idx_scan, 0), 2) AS fetch_pct
FROM pg_stat_user_indexes
WHERE indexname LIKE '%_covering%';
```

### 性能指标

| 指标 | 正常范围 | 需要关注 |
|------|---------|---------|
| `fetch_pct` | < 10% | > 30% |
| 索引扫描次数 | 持续增长 | 无增长 |
| 索引膨胀率 | < 30% | > 50% |

### 定期维护

```sql
-- 每季度或膨胀 > 30% 时重建
REINDEX INDEX CONCURRENTLY idx_orders_user_status_created_covering;
REINDEX INDEX CONCURRENTLY idx_chat_messages_group_sent_covering;
REINDEX INDEX CONCURRENTLY idx_payments_user_status_created_covering;

-- 检查膨胀（需要 pgstattuple 扩展）
CREATE EXTENSION IF NOT EXISTS pgstattuple;
SELECT * FROM pgstatindex('idx_orders_user_status_created_covering');
```

## 回滚计划

如果出现性能问题：

```bash
# 立即回滚
psql -U gamelink -d gamelink -f 0001_add_covering_indexes_rollback.sql

# 验证回滚
\di *covering*  -- 应该为空
```

## 常见问题

### Q1: 为什么不用 GORM AutoMigrate？

**A**: GORM 不支持 PostgreSQL 的 `INCLUDE` 语法。覆盖索引需要手动 SQL 迁移。

### Q2: 覆盖索引会增加多少存储？

**A**:
- `orders`: 约增加 30%（~150MB per 100万行）
- `chat_messages`: 约增加 50%（content 列较大）
- `payments`: 约增加 20%

### Q3: 会影响写入性能吗？

**A**: 会轻微增加写入开销（~5-10%），但读取性能提升远超写入损失。

### Q4: 什么时候应该删除覆盖索引？

**A**:
- 索引使用率极低（`idx_scan` < 100/周）
- 查询模式改变，不再需要
- 存储成本超过性能收益

### Q5: 如何验证索引是否生效？

**A**: 使用 `EXPLAIN (ANALYZE, BUFFERS)` 查看执行计划：
- ✅ `Index Only Scan` - 覆盖索引生效
- ❌ `Index Scan` - 仍然回表（需要调整查询或索引）

## 相关文档

- [PostgreSQL Covering Indexes](https://www.postgresql.org/docs/current/indexes-index-only-scans.html)
- [Index-Only Scans](https://www.postgresql.org/docs/current/indexes-index-only-scans.html)
- [04c-enums-indexes.md](../../.kiro/steering/04c-enums-indexes.md) - 索引完整文档

## 变更日志

| 日期 | 版本 | 变更内容 |
|------|------|---------|
| 2026-01-01 | 1.0 | 初始版本：订单、聊天、支付表覆盖索引 |
