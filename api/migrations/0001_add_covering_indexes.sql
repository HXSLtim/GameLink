-- ============================================================================
-- GameLink Covering Index Migration
-- ============================================================================
-- Purpose: Add covering indexes to optimize common queries and reduce table
--          access (heap fetch) by including frequently accessed columns.
--
-- Author: DBA Team
-- Date: 2026-01-01
-- Version: 1.0
--
-- Notes:
-- - Uses PostgreSQL INCLUDE clause (requires PostgreSQL 11+)
-- - Run with CONCURRENTLY to avoid blocking production traffic
-- - Covering indexes add storage overhead but improve query performance
--
-- Migration Commands:
--   psql -U gamelink -d gamelink -f 0001_add_covering_indexes.sql
--
-- Rollback:
--   See 0001_add_covering_indexes_rollback.sql
-- ============================================================================

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';

-- ----------------------------------------------------------------------------
-- 1. Order Table Covering Index
-- ----------------------------------------------------------------------------
-- Query Pattern:
--   SELECT id, player_id, total_price_cents
--   FROM orders
--   WHERE user_id = ? AND status IN (?)
--   ORDER BY created_at DESC LIMIT 20;
--
-- Benefit:
--   - Covers index-only scan for user order list queries
--   - Includes player_id (foreign key for join)
--   - Includes total_price_cents (display in list)
--   - Eliminates heap fetch for common list operations
-- ----------------------------------------------------------------------------

-- Drop old partial index if exists (safe to run, ignores if not exists)
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_user_status_created;

-- Create new covering index
CREATE INDEX CONCURRENTLY idx_orders_user_status_created_covering
ON orders (user_id, status, created_at DESC)
INCLUDE (id, player_id, total_price_cents, commission_cents, player_income_cents);

-- Comment the index
COMMENT ON INDEX idx_orders_user_status_created_covering IS
'Covering index for user order list queries. Includes (user_id, status, created_at) with INCLUDE columns (id, player_id, total_price_cents, commission_cents, player_income_cents). Optimizes: SELECT id, player_id, total_price_cents FROM orders WHERE user_id = ? AND status IN (?) ORDER BY created_at DESC.';

-- ----------------------------------------------------------------------------
-- 2. Chat Messages Table Covering Index
-- ----------------------------------------------------------------------------
-- Query Pattern:
--   SELECT id, content, sender_id, created_at
--   FROM chat_messages
--   WHERE group_id = ?
--   ORDER BY created_at DESC
--   LIMIT 50;
--
-- Benefit:
--   - Covers index-only scan for message list queries
--   - Includes content (message text, frequently accessed)
--   - Includes sender_id (for displaying sender info)
--   - Eliminates heap fetch for message history
-- ----------------------------------------------------------------------------

-- Drop old partial index if exists
DROP INDEX CONCURRENTLY IF EXISTS idx_chat_messages_group_created;

-- Create new covering index
CREATE INDEX CONCURRENTLY idx_chat_messages_group_sent_covering
ON chat_messages (group_id, created_at DESC)
INCLUDE (id, content, sender_id, message_type, audit_status);

-- Comment the index
COMMENT ON INDEX idx_chat_messages_group_sent_covering IS
'Covering index for chat message list queries. Includes (group_id, created_at) with INCLUDE columns (id, content, sender_id, message_type, audit_status). Optimizes: SELECT id, content, sender_id FROM chat_messages WHERE group_id = ? ORDER BY created_at DESC LIMIT 50.';

-- ----------------------------------------------------------------------------
-- 3. Payment Table Covering Index (Bonus Optimization)
-- ----------------------------------------------------------------------------
-- Query Pattern:
--   SELECT id, amount_cents, status, payment_method
--   FROM payments
--   WHERE user_id = ? AND status IN (?)
--   ORDER BY created_at DESC;
--
-- Benefit:
--   - Covers payment history queries
--   - Includes amount and payment method for display
-- ----------------------------------------------------------------------------

CREATE INDEX CONCURRENTLY idx_payments_user_status_created_covering
ON payments (user_id, status, created_at DESC)
INCLUDE (id, amount_cents, payment_method, provider_trade_no);

COMMENT ON INDEX idx_payments_user_status_created_covering IS
'Covering index for payment history queries. Includes (user_id, status, created_at) with INCLUDE columns (id, amount_cents, payment_method, provider_trade_no). Optimizes payment list display.';

-- ----------------------------------------------------------------------------
-- 4. Verification and Analysis
-- ----------------------------------------------------------------------------

-- Display index sizes for comparison
SELECT
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
WHERE tablename IN ('orders', 'chat_messages', 'payments')
ORDER BY tablename, indexname;

-- Display index usage statistics (run after some traffic)
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan AS index_scans,
    idx_tup_read AS tuples_read,
    idx_tup_fetch AS tuples_fetched,
    pg_size_pretty(pg_relation_size(indexrelid)) AS size
FROM pg_stat_user_indexes
WHERE tablename IN ('orders', 'chat_messages', 'payments')
ORDER BY idx_scan DESC;

-- ============================================================================
-- Migration Complete
-- ============================================================================
-- Next Steps:
-- 1. Run EXPLAIN ANALYZE on target queries to verify index-only scans
-- 2. Monitor pg_stat_user_indexes for index usage
-- 3. Set up alerts for index bloat (pgstattuple extension)
-- 4. Schedule regular index maintenance (REINDEX CONCURRENTLY)
-- ============================================================================
