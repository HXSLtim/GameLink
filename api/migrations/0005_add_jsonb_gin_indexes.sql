-- ============================================================================
-- GameLink JSONB GIN Index Migration
-- ============================================================================
-- Purpose: Add GIN indexes for JSONB columns to optimize JSON query performance
--          and enable efficient JSON path queries using @>, ?, ?&, ?| operators.
--
-- Author: Database Team
-- Date: 2026-03-02
-- Version: 1.0
--
-- Performance Impact:
-- - Enables index-only scans for JSON containment queries (@>)
-- - Speeds up JSON key existence checks (?, ?&, ?|)
-- - Reduces query time from O(n) table scan to O(log n) index lookup
-- - Expected improvement: 10-100x faster for JSON queries
--
-- Storage Impact:
-- - GIN indexes are larger than B-tree (typically 2-3x column size)
-- - Monitor disk usage after creation
--
-- Migration Commands:
--   psql -U gamelink -d gamelink -f 0005_add_jsonb_gin_indexes.sql
--
-- Rollback:
--   See 0005_add_jsonb_gin_indexes_rollback.sql
-- ============================================================================

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';

-- ============================================================================
-- SECTION 1: Base Model ExtJSON Field (All Tables)
-- ============================================================================
-- The ext_json field exists in all tables via Base model.
-- Priority: HIGH - affects all entities for future extensibility.
-- ============================================================================

-- Users table
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_ext_json_gin
ON users USING GIN (ext_json jsonb_path_ops);

COMMENT ON INDEX idx_users_ext_json_gin IS
'GIN index for users.ext_json JSONB queries. Optimizes @> containment queries.';

-- Players table
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_players_ext_json_gin
ON players USING GIN (ext_json jsonb_path_ops);

COMMENT ON INDEX idx_players_ext_json_gin IS
'GIN index for players.ext_json JSONB queries.';

-- Orders table
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_ext_json_gin
ON orders USING GIN (ext_json jsonb_path_ops);

COMMENT ON INDEX idx_orders_ext_json_gin IS
'GIN index for orders.ext_json JSONB queries.';

-- Payments table
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_payments_ext_json_gin
ON payments USING GIN (ext_json jsonb_path_ops);

COMMENT ON INDEX idx_payments_ext_json_gin IS
'GIN index for payments.ext_json JSONB queries.';

-- Chat groups table
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_groups_ext_json_gin
ON chat_groups USING GIN (ext_json jsonb_path_ops);

COMMENT ON INDEX idx_chat_groups_ext_json_gin IS
'GIN index for chat_groups.ext_json JSONB queries.';

-- Chat messages table
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_messages_ext_json_gin
ON chat_messages USING GIN (ext_json jsonb_path_ops);

COMMENT ON INDEX idx_chat_messages_ext_json_gin IS
'GIN index for chat_messages.ext_json JSONB queries.';

-- Reviews table
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reviews_ext_json_gin
ON reviews USING GIN (ext_json jsonb_path_ops);

COMMENT ON INDEX idx_reviews_ext_json_gin IS
'GIN index for reviews.ext_json JSONB queries.';

-- ============================================================================
-- SECTION 2: Business-Specific JSON Columns
-- ============================================================================

-- ----------------------------------------------------------------------------
-- 2.1 Orders: order_config (JSON configuration)
-- ----------------------------------------------------------------------------
-- Query pattern: WHERE order_config @> '{"key": "value"}'
-- Use case: Filter orders by custom configuration attributes
-- ----------------------------------------------------------------------------

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_config_gin
ON orders USING GIN (
    CAST(order_config AS jsonb) jsonb_path_ops
);

COMMENT ON INDEX idx_orders_config_gin IS
'GIN index for orders.order_config JSON queries. Enables fast filtering by order configuration.';

-- ----------------------------------------------------------------------------
-- 2.2 Chat Groups: settings (JSON settings)
-- ----------------------------------------------------------------------------
-- Query pattern: WHERE settings @> '{"notifications": true}'
-- Use case: Filter chat groups by settings
-- ----------------------------------------------------------------------------

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_groups_settings_gin
ON chat_groups USING GIN (
    CAST(settings AS jsonb) jsonb_path_ops
);

COMMENT ON INDEX idx_chat_groups_settings_gin IS
'GIN index for chat_groups.settings JSON queries.';

-- ----------------------------------------------------------------------------
-- 2.3 Chat Messages: metadata (JSON metadata)
-- ----------------------------------------------------------------------------
-- Query pattern: WHERE metadata @> '{"type": "system"}'
-- Use case: Filter messages by metadata attributes
-- ----------------------------------------------------------------------------

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_messages_metadata_gin
ON chat_messages USING GIN (
    CAST(metadata AS jsonb) jsonb_path_ops
);

COMMENT ON INDEX idx_chat_messages_metadata_gin IS
'GIN index for chat_messages.metadata JSON queries.';

-- ----------------------------------------------------------------------------
-- 2.4 Payments: provider_raw (JSON provider response)
-- ----------------------------------------------------------------------------
-- Query pattern: WHERE provider_raw @> '{"status": "success"}'
-- Use case: Query payments by provider response details
-- ----------------------------------------------------------------------------

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_payments_provider_raw_gin
ON payments USING GIN (provider_raw jsonb_path_ops);

COMMENT ON INDEX idx_payments_provider_raw_gin IS
'GIN index for payments.provider_raw JSONB queries. Optimizes provider response filtering.';

-- ----------------------------------------------------------------------------
-- 2.5 Operation Logs: metadata_json (JSON metadata)
-- ----------------------------------------------------------------------------
-- Query pattern: WHERE metadata_json @> '{"action": "login"}'
-- Use case: Filter operation logs by metadata
-- Priority: HIGH - operation_logs is high-volume table
-- ----------------------------------------------------------------------------

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_operation_logs_metadata_gin
ON operation_logs USING GIN (metadata_json jsonb_path_ops);

COMMENT ON INDEX idx_operation_logs_metadata_gin IS
'GIN index for operation_logs.metadata_json JSONB queries. Critical for log analysis.';

-- ----------------------------------------------------------------------------
-- 2.6 Financial Reports: report_data (JSON report data)
-- ----------------------------------------------------------------------------
-- Query pattern: WHERE report_data @> '{"period": "2024-01"}'
-- Use case: Filter financial reports by data attributes
-- ----------------------------------------------------------------------------

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_financial_reports_data_gin
ON financial_reports USING GIN (report_data jsonb_path_ops);

COMMENT ON INDEX idx_financial_reports_data_gin IS
'GIN index for financial_reports.report_data JSONB queries.';

-- ----------------------------------------------------------------------------
-- 2.7 User Settings: notifications, privacy (JSON settings)
-- ----------------------------------------------------------------------------
-- Query pattern: WHERE notifications @> '{"email": true}'
-- Use case: Filter users by notification preferences
-- ----------------------------------------------------------------------------

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_settings_notifications_gin
ON user_settings USING GIN (
    CAST(notifications AS jsonb) jsonb_path_ops
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_settings_privacy_gin
ON user_settings USING GIN (
    CAST(privacy AS jsonb) jsonb_path_ops
);

COMMENT ON INDEX idx_user_settings_notifications_gin IS
'GIN index for user_settings.notifications JSON queries.';

COMMENT ON INDEX idx_user_settings_privacy_gin IS
'GIN index for user_settings.privacy JSON queries.';

-- ----------------------------------------------------------------------------
-- 2.8 VIP Levels: benefits (JSON benefits)
-- ----------------------------------------------------------------------------
-- Query pattern: WHERE benefits @> '{"discount": 0.2}'
-- Use case: Query VIP levels by benefit attributes
-- ----------------------------------------------------------------------------

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_vip_levels_benefits_gin
ON vip_levels USING GIN (
    CAST(benefits AS jsonb) jsonb_path_ops
);

COMMENT ON INDEX idx_vip_levels_benefits_gin IS
'GIN index for vip_levels.benefits JSON queries.';

-- ============================================================================
-- SECTION 3: Verification and Monitoring
-- ============================================================================

-- Display all new GIN indexes
SELECT
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size,
    idx_scan AS scans,
    idx_tup_read AS tuples_read
FROM pg_stat_user_indexes
WHERE indexname LIKE '%_gin'
ORDER BY tablename, indexname;

-- Check index bloat (requires pgstattuple extension)
-- Uncomment if pgstattuple is installed:
-- SELECT
--     schemaname,
--     tablename,
--     indexname,
--     pg_size_pretty(pg_relation_size(indexrelid)) AS index_size,
--     round(100 * pgstatindex(indexrelid)::numeric / pg_relation_size(indexrelid), 2) AS bloat_pct
-- FROM pg_stat_user_indexes
-- WHERE indexname LIKE '%_gin';

-- ============================================================================
-- Migration Complete
-- ============================================================================
-- Next Steps:
-- 1. Monitor index usage with pg_stat_user_indexes
-- 2. Run EXPLAIN ANALYZE on JSON queries to verify index usage
-- 3. Set up alerts for index bloat (>30% bloat requires REINDEX)
-- 4. Schedule regular VACUUM ANALYZE for statistics updates
--
-- Example Query to Verify Index Usage:
-- EXPLAIN (ANALYZE, BUFFERS) SELECT * FROM orders WHERE order_config @> '{"priority": "high"}';
-- Look for "Index Scan using idx_orders_config_gin" in the plan
-- ============================================================================
