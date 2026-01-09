-- ============================================================================
-- GameLink Foreign Key Index Migration
-- ============================================================================
-- Purpose: Add indexes on foreign key columns to optimize JOIN operations
--          and improve query performance for related data lookups.
--
-- Author: Performance Team
-- Date: 2026-01-09
-- Version: 1.0
--
-- Notes:
-- - Foreign key columns without indexes cause slow JOINs
-- - Run with CONCURRENTLY to avoid blocking production traffic
-- - These indexes are essential for N+1 query optimization
--
-- Migration Commands:
--   psql -U gamelink -d gamelink -f 0002_add_foreign_key_indexes.sql
-- ============================================================================

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';

-- ----------------------------------------------------------------------------
-- 1. Users Table - Foreign Key Indexes
-- ----------------------------------------------------------------------------

-- Index for referrer_id (used in referral queries)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_referrer_id
ON users (referrer_id)
WHERE referrer_id IS NOT NULL;

COMMENT ON INDEX idx_users_referrer_id IS
'Index for referrer lookups in referral system queries.';

-- ----------------------------------------------------------------------------
-- 2. Chat Messages Table - Foreign Key Indexes
-- ----------------------------------------------------------------------------

-- Index for sender_id (used in message sender lookups)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_messages_sender_id
ON chat_messages (sender_id);

COMMENT ON INDEX idx_chat_messages_sender_id IS
'Index for sender lookups when displaying chat messages.';

-- Index for moderated_by (used in moderation queries)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_messages_moderated_by
ON chat_messages (moderated_by)
WHERE moderated_by IS NOT NULL;

COMMENT ON INDEX idx_chat_messages_moderated_by IS
'Index for moderator activity tracking.';

-- ----------------------------------------------------------------------------
-- 3. Payments Table - Foreign Key Indexes
-- ----------------------------------------------------------------------------

-- Index for order_id (used in order payment lookups)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_payments_order_id
ON payments (order_id)
WHERE order_id IS NOT NULL;

COMMENT ON INDEX idx_payments_order_id IS
'Index for order payment lookups.';

-- ----------------------------------------------------------------------------
-- 4. Reviews Table - Foreign Key Indexes
-- ----------------------------------------------------------------------------

-- Index for user_id (reviewer)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reviews_user_id
ON reviews (user_id);

COMMENT ON INDEX idx_reviews_user_id IS
'Index for user review history queries.';

-- Index for player_id (reviewed player)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reviews_player_id
ON reviews (player_id);

COMMENT ON INDEX idx_reviews_player_id IS
'Index for player review aggregation queries.';

-- Index for order_id
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reviews_order_id
ON reviews (order_id);

COMMENT ON INDEX idx_reviews_order_id IS
'Index for order review lookups.';

-- Composite index for player rating queries
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reviews_player_status_rating
ON reviews (player_id, status, rating)
WHERE status = 'approved';

COMMENT ON INDEX idx_reviews_player_status_rating IS
'Composite index for player rating calculation queries.';

-- ----------------------------------------------------------------------------
-- 5. Orders Table - Additional Foreign Key Indexes
-- ----------------------------------------------------------------------------

-- Index for player_id (陪玩师订单查询)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_player_id
ON orders (player_id)
WHERE player_id IS NOT NULL;

COMMENT ON INDEX idx_orders_player_id IS
'Index for player order list queries.';

-- Index for game_id (游戏订单统计)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_game_id
ON orders (game_id)
WHERE game_id IS NOT NULL;

COMMENT ON INDEX idx_orders_game_id IS
'Index for game order statistics.';

-- ----------------------------------------------------------------------------
-- 6. Wallets Table - Foreign Key Indexes
-- ----------------------------------------------------------------------------

-- Index for user_id (钱包查询)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_wallets_user_id
ON wallets (user_id);

COMMENT ON INDEX idx_wallets_user_id IS
'Index for user wallet lookups.';

-- ----------------------------------------------------------------------------
-- 7. Notifications Table - Foreign Key Indexes
-- ----------------------------------------------------------------------------

-- Index for user_id (用户通知查询)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_notifications_user_id_read
ON notifications (user_id, is_read, created_at DESC);

COMMENT ON INDEX idx_notifications_user_id_read IS
'Composite index for user notification list with read status filter.';

-- ----------------------------------------------------------------------------
-- 8. Verification and Analysis
-- ----------------------------------------------------------------------------

-- Display all new indexes
SELECT
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
WHERE indexname LIKE 'idx_%'
  AND tablename IN ('users', 'chat_messages', 'payments', 'reviews', 'orders', 'wallets', 'notifications')
ORDER BY tablename, indexname;

-- ============================================================================
-- Migration Complete
-- ============================================================================
